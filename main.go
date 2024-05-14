package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/BeeTimeClock/BeeTimeClock-Server/auth"
	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/database"
	"github.com/BeeTimeClock/BeeTimeClock-Server/handler"
	"github.com/BeeTimeClock/BeeTimeClock-Server/middleware"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
)

var (
	GitCommit string
)

const MIGRATION_HOMEOFFICE_GOING = "HOMEOFFICE_GOING"

func main() {
	env := core.NewEnvironment()

	db := database.NewDatabaseManager("beetc")
	env.DatabaseManager = db

	userRepo := repository.NewUser(env)
	err := userRepo.Migrate()
	if err != nil {
		panic(err)
	}

	timestampRepo := repository.NewTimestamp(env)
	err = timestampRepo.Migrate()
	if err != nil {
		panic(err)
	}

	fuelRepo := repository.NewFuel(env)
	err = fuelRepo.Migrate()
	if err != nil {
		panic(err)
	}

	absenceRepo := repository.NewAbsence(env)
	err = absenceRepo.Migrate()
	if err != nil {
		panic(err)
	}

	migrationRepo := repository.NewMigration(env)
	err = migrationRepo.Migrate()
	if err != nil {
		panic(err)
	}

	settingsRepo := repository.NewSettings(env)
	err = settingsRepo.Migrate()
	if err != nil {
		panic(err)
	}

	userHandler := handler.NewUser(env, userRepo)
	timestampHandler := handler.NewTimestamp(env, userRepo, timestampRepo, absenceRepo, settingsRepo)
	fuelHandler := handler.NewFuel(env, userRepo, fuelRepo)
	absenceHandler := handler.NewAbsence(env, userRepo, absenceRepo)
	migrationHandler := handler.NewMigration(env, migrationRepo)
	administrationHandler := handler.NewAdministration(env, settingsRepo)

	authProvider := auth.NewAuthProvider(env, userRepo)
	if err != nil {
		panic(err)
	}

	_, err = migrationRepo.MigrationFindByTitle(MIGRATION_HOMEOFFICE_GOING)
	homeofficeGoingMigrationExists := true
	if err != nil {
		if err == repository.ErrMigrationNotFound {
			homeofficeGoingMigrationExists = false
		} else {
			panic(err)
		}
	}

	if !homeofficeGoingMigrationExists {
		log.Println("Migration: HOMEOFFICE_GOING started")
		timestamps, err := timestampRepo.FindAll()
		if err != nil {
			panic(err)
		}

		for _, timestamp := range timestamps {
			timestamp.IsHomeofficeGoing = timestamp.IsHomeoffice
			err = timestampRepo.Update(&timestamp)

			if err != nil {
				homeofficeGoingMigration := model.Migration{
					Title:      MIGRATION_HOMEOFFICE_GOING,
					Result:     err.Error(),
					FinishedAt: time.Now(),
					Success:    true,
				}
				migrationRepo.MigrationInsert(&homeofficeGoingMigration)

				panic(err)
			}
		}

		homeofficeGoingMigration := model.Migration{
			Title:      MIGRATION_HOMEOFFICE_GOING,
			Result:     fmt.Sprintf("%d timestamps were migrated", len(timestamps)),
			FinishedAt: time.Now(),
			Success:    true,
		}
		migrationRepo.MigrationInsert(&homeofficeGoingMigration)
		log.Println("Migration: HOMEOFFICE_GOING finished")
	} else {
		log.Println("Migration: HOMEOFFICE_GOING already finished")
	}

	r := gin.Default()
	r.Use(middleware.AcceptCors)

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Next()
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, "/ui/")
	})

	uiFSSub, _ := fs.Sub(uiFS, "ui/dist/spa")
	r.StaticFS("/ui/", &uiWrapper{FileSystem: http.FS(uiFSSub)})

	v1 := r.Group("api/v1")
	{
		v1.GET("auth", authProvider.Auth)
		v1.GET("auth/providers", authProvider.AuthProviders)
		v1.GET("auth/microsoft", authProvider.MicrosoftAuthSettings)

		v1.GET("status", func(c *gin.Context) {
			commit := GitCommit
			if commit == "" {
				commit = "dirty"
			}

			c.JSON(http.StatusOK, model.NewSuccessResponse(gin.H{
				"Commit": commit,
			}))
		})

		v1.Use(authProvider.AuthRequired)
		{
			administration := v1.Group("administration")
			{
				administration.Use(auth.AdministratorAccessRequired)
				administrationUser := administration.Group("user")
				{
					administrationUser.GET("", userHandler.AdministrationUserGetAll)
					administrationUser.POST("", userHandler.AdministrationUserCreate)
					administrationUser.PUT(":userID", userHandler.AdministrationUserUpdate)
					administrationUser.GET(":userID", userHandler.AdministrationUserGetByUserID)
					administrationUser.DELETE(":userID", userHandler.AdministrationUserDelete)

					administrationUser.GET(":userID/absence/year/:year/summary", absenceHandler.AbsenceQueryUserSummaryYear)
					administrationUser.GET(":userID/absence/year/:year", absenceHandler.AbsenceQueryUserYear)
					administrationUser.GET(":userID/absence/years", absenceHandler.AbsenceQueryUserYears)

					administrationUser.GET(":userID/timestamp/year/:year/month/:month/grouped", timestampHandler.TimestampUserQueryMonthGrouped)
					administrationUser.GET(":userID/timestamp/year/:year/month/:month/overtime", timestampHandler.TimestampUserQueryMonthOvertime)
					administrationUser.GET(":userID/timestamp/months", timestampHandler.TimestampUserQueryMonths)
				}
				administrationMigrations := administration.Group("migration")
				{
					administrationMigrations.GET("", migrationHandler.AdministrationMigrationGetAll)
				}
				administrationSettings := administration.Group("settings")
				{
					administrationSettings.GET("", administrationHandler.AdministrationGetSettings)
					administrationSettings.PUT("", administrationHandler.AdministrationUpdateSettings)
				}
			}

			timestamp := v1.Group("timestamp")
			{
				timestamp.GET("", timestampHandler.TimestampGetAll)
				timestamp.GET("query/last", timestampHandler.TimestampQueryLast)
				timestamp.GET("query/current_month/grouped", timestampHandler.TimestampCurrentUserQueryCurrentMonthGrouped)
				timestamp.GET("query/current_month/overtime", timestampHandler.TimestampCurrentUserQueryCurrentMonthOvertime)
				timestamp.GET("query/year/:year/month/:month/grouped", timestampHandler.TimestampQueryMonthGrouped)
				timestamp.GET("query/year/:year/month/:month/overtime", timestampHandler.TimestampQueryMonthOvertime)
				timestamp.GET("query/timestamp/months", timestampHandler.TimestampQueryMonths)
				timestamp.POST("action/checkin", timestampHandler.TimestampActionCheckIn)
				timestamp.POST("action/checkout", timestampHandler.TimestampActionCheckOut)
				timestamp.POST(":timestampID/correction", timestampHandler.TimestampCorrectionCreate)
				timestamp.POST("", timestampHandler.TimestampCreate)
				timestamp.GET("overtime", timestampHandler.TimestampOvertime)
			}

			fuel := v1.Group("fuel")
			{
				fuel.GET("", fuelHandler.FuelGetAll)
				fuel.GET(":fuelID", fuelHandler.FuelGet)
				fuel.PUT(":fuelID", fuelHandler.FuelUpdate)
				fuel.POST("action/prepare", fuelHandler.FuelActionPrepare)
			}

			absence := v1.Group("absence")
			{
				absence.GET("", absenceHandler.AbsenceGetAll)
				absence.POST("", absenceHandler.AbsenceCreate)
				absence.DELETE(":id", absenceHandler.AbsenceDelete)
				absence.GET("query/me/summary", absenceHandler.AbsenceQueryCurrentUserSummary)
				absence.GET("query/users/summary", absenceHandler.AbsenceQueryUsersSummary)
				absence.GET("reasons", absenceHandler.AbsenceReasonsGetAll)
			}

			user := v1.Group("user")
			{
				user.GET("me", userHandler.CurrentUserGet)
				user.GET("me/apikey", userHandler.CurrentUserApikeyGet)
				user.POST("me/apikey", userHandler.CurrentUserApikeyCreate)
			}
		}
	}

	r.Run()
}

func importHolidays(env *core.Environment, absenceRepo *repository.Absence, year int) error {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://feiertage-api.de/api/?jahr=%d&nur_land=NI", year), nil)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	result := make(map[string]model.HolidayImport)

	body, _ := io.ReadAll(response.Body)

	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	for name, info := range result {
		date, err := info.GetDate()
		if err != nil {
			return err
		}

		exists, err := absenceRepo.HolidayIsByDate(date)
		if err != nil {
			return err
		}

		if !exists {
			err = absenceRepo.HolidayInsert(&model.Holiday{
				Name: name,
				Date: date,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type uiWrapper struct {
	FileSystem http.FileSystem
}

func (w *uiWrapper) Open(name string) (http.File, error) {
	// return file if it exists
	file, err := w.FileSystem.Open(name)
	if err == nil {
		return file, nil
	}

	// redirect non-existing files to index.html
	// required for spa ui to work correctly
	if errors.Is(err, fs.ErrNotExist) {
		file, err := w.FileSystem.Open("index.html")
		return file, err
	}

	return nil, err
}
