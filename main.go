package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"

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

	userHandler := handler.NewUser(env, userRepo)
	timestampHandler := handler.NewTimestamp(env, userRepo, timestampRepo, absenceRepo)
	fuelHandler := handler.NewFuel(env, userRepo, fuelRepo)
	absenceHandler := handler.NewAbsence(env, userRepo, absenceRepo)

	authProvider := auth.NewAuthProvider(env, userRepo)
	if err != nil {
		panic(err)
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
	r.StaticFS("ui", http.FS(uiFSSub))

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

					administrationUser.GET(":userID/timestamp/year/:year/month/:month/grouped", timestampHandler.TimestampUserQueryMonthGrouped)
					administrationUser.GET(":userID/timestamp/year/:year/month/:month/overtime", timestampHandler.TimestampUserQueryMonthOvertime)
					administrationUser.GET(":userID/timestamp/months", timestampHandler.TimestampUserQueryMonths)
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
