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
	"github.com/BeeTimeClock/BeeTimeClock-Server/microsoft"
	"github.com/BeeTimeClock/BeeTimeClock-Server/middleware"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/BeeTimeClock/BeeTimeClock-Server/worker"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	GitCommit string
)

const (
	MIGRATION_HOMEOFFICE_GOING        = "HOMEOFFICE_GOING"
	MIGRATION_EXTERNAL_CALENDAR       = "EXTERNAL_CALENDAR"
	MIGRATION_EXTERNAL_CALENDAR_MULTI = "EXTERNAL_CALENDAR_MULTI"
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

	teamRepo := repository.NewTeam(env)
	err = teamRepo.Migrate()
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

	externalWorkRepo := repository.NewExternalWork(env)
	err = externalWorkRepo.Migrate()
	if err != nil {
		panic(err)
	}

	overtimeRepo := repository.NewOvertime(env)
	err = overtimeRepo.Migrate()
	if err != nil {
		panic(err)
	}

	holiday := repository.NewHoliday(env)
	err = holiday.Migrate()
	if err != nil {
		panic(err)
	}

	timestampWorker := worker.NewTimestamp(env, userRepo, externalWorkRepo, timestampRepo, holiday)
	overtimeWorker := worker.NewOvertime(env, userRepo, externalWorkRepo, timestampRepo, holiday, overtimeRepo, timestampWorker, absenceRepo)

	userHandler := handler.NewUser(env, userRepo, teamRepo)
	timestampHandler := handler.NewTimestamp(env, userRepo, timestampRepo, absenceRepo, settingsRepo, holiday, timestampWorker)
	fuelHandler := handler.NewFuel(env, userRepo, fuelRepo)
	absenceHandler := handler.NewAbsence(env, userRepo, absenceRepo)
	migrationHandler := handler.NewMigration(env, migrationRepo)
	administrationHandler := handler.NewAdministration(env, settingsRepo, absenceRepo)
	externalWorkHandler := handler.NewExternalWork(env, userRepo, externalWorkRepo, holiday)
	overtimeHandler := handler.NewOvertime(env, userRepo, overtimeRepo, overtimeWorker)

	authProvider := auth.NewAuthProvider(env, userRepo)

	go importHolidays(holiday)
	go overtimeWorker.CalculateMissingMonths()

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

	err = migrateExternalCalendar(migrationRepo, absenceRepo)
	if err != nil {
		panic(err)
	}

	err = migrateExternalCalendarMulti(migrationRepo, absenceRepo)
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
	r.StaticFS("/ui/", &uiWrapper{FileSystem: http.FS(uiFSSub)})

	v1 := r.Group("api/v1")
	{
		v1.GET("logo", administrationHandler.GetLogo)
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
				administrationTeam := administration.Group("team")
				{
					administrationTeam.GET("", userHandler.AdministrationTeamGetAll)
					administrationTeam.POST("", userHandler.AdministrationTeamCreate)
					administrationTeam.PUT(":teamID", userHandler.AdministrationTeamUpdate)
					administrationTeam.GET(":teamID", userHandler.AdministrationTeamGetByID)
					administrationTeam.DELETE(":teamID", userHandler.AdministrationTeamDelete)
					administrationTeam.GET(":teamID/member", userHandler.AdministrationTeamMemberGetByTeamID)
					administrationTeam.POST(":teamID/member", userHandler.AdministrationTeamMemberCreate)
					administrationTeam.DELETE(":teamID/member/:teamMemberID", userHandler.AdministrationTeamMemberDelete)
				}
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
				administrationAbsence := administration.Group("absence")
				{
					administrationAbsence.GET("reasons", absenceHandler.AbsenceReasonsGetAll)
					administrationAbsence.POST("reasons", absenceHandler.AdministrationAbsenceReasonCreate)
					administrationAbsence.PUT("reasons/:absenceReasonID", absenceHandler.AdministrationAbsenceReasonUpdate)
					administrationAbsence.DELETE("reasons/:absenceReasonID", absenceHandler.AdministrationAbsenceReasonDelete)
				}

				administrationExternalWork := administration.Group("external_work")
				{
					administrationExternalWork.GET("compensation", externalWorkHandler.AdministrationExternalWorkCompensationGetAll)
					administrationExternalWork.POST("compensation", externalWorkHandler.AdministrationExternalWorkCompensationCreate)
					administrationExternalWork.PUT("compensation/:externalWorkCompensationId", externalWorkHandler.AdministrationExternalWorkCompensationUpdate)
				}

				administrationMigrations := administration.Group("migration")
				{
					administrationMigrations.GET("", migrationHandler.AdministrationMigrationGetAll)
				}
				administrationSettings := administration.Group("settings")
				{
					administrationSettings.GET("", administrationHandler.AdministrationGetSettings)
					administrationSettings.PUT("", administrationHandler.AdministrationUpdateSettings)
					administrationSettings.POST("logo", administrationHandler.AdministrationUploadLogo)
				}
				administrationNotify := administration.Group("notify")
				{
					administrationNotify.POST("absence/week", administrationHandler.AdministrationNotifyAbsenceWeek)
				}
			}

			timestamp := v1.Group("timestamp")
			{
				timestamp.GET("", timestampHandler.TimestampGetAll)
				timestamp.GET("query/last", timestampHandler.TimestampQueryLast)
				timestamp.GET("query/suspicious", timestampHandler.TimestampQuerySuspicious)
				timestamp.GET("query/suspicious/count", timestampHandler.TimestampQuerySuspiciousCount)
				timestamp.GET("query/current_month/grouped", timestampHandler.TimestampCurrentUserQueryCurrentMonthGrouped)
				timestamp.GET("query/current_month/overtime", timestampHandler.TimestampCurrentUserQueryCurrentMonthOvertime)
				timestamp.GET("query/year/:year/month/:month/grouped", timestampHandler.TimestampQueryMonthGrouped)
				timestamp.GET("query/year/:year/month/:month/overtime", timestampHandler.TimestampQueryMonthOvertime)
				timestamp.GET("query/timestamp/months", timestampHandler.TimestampQueryMonths)
				timestamp.POST("action/checkin", timestampHandler.TimestampActionCheckIn)
				timestamp.POST("action/checkout", timestampHandler.TimestampActionCheckOut)
				timestamp.POST(":timestampID/correction", timestampHandler.TimestampCorrectionCreate)
				timestamp.POST("", timestampHandler.TimestampCreate)
			}

			overtime := v1.Group("overtime")
			{
				overtime.GET("", overtimeHandler.OvertimeGetAll)
				overtime.GET("total", overtimeHandler.OvertimeTotal)
				overtime.POST("action/calculate/:year/:month", overtimeHandler.OvertimeCurrentUserCalculateMonth)
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
				absence.GET("query/users/summary/current_year", absenceHandler.AbsenceQueryUsersSummaryCurrentYear)
				absence.GET("query/users/summary/current_week", absenceHandler.AbsenceQueryUsersSummaryCurrentWeek)
				absence.GET("reasons", absenceHandler.AbsenceReasonsGetAll)
			}

			externalWork := v1.Group("external_work")
			{
				externalWork.GET("", externalWorkHandler.ExternalWorkGetAll)
				externalWork.POST("", externalWorkHandler.ExternalWorkCreate)
				externalWork.GET("invoiced", externalWorkHandler.ExternalWorkGetInvoiced)
				externalWork.GET("action/export/pdf", externalWorkHandler.ExternalWorkExportPdf)
				externalWork.GET("action/export/pdf/:invoiceIdentifier", externalWorkHandler.ExternalWorkDownloadPdf)

				externalWorkDetail := externalWork.Group(":externalWorkId")
				{
					externalWorkDetail.GET("", externalWorkHandler.ExternalWorkGetById)
					externalWorkDetail.POST("expanse", externalWorkHandler.ExternalWorkExpanseCreate)
					externalWorkDetail.PUT("expanse/:externalWorkExpanseId", externalWorkHandler.ExternalWorkExpanseUpdate)
					externalWorkDetail.POST("action/submit", externalWorkHandler.ExternalWorkSubmit)
				}
			}

			user := v1.Group("user")
			{
				user.GET("me", userHandler.CurrentUserGet)
				user.PUT("me", userHandler.CurrentUserUpdate)
				user.GET("me/apikey", userHandler.CurrentUserApikeyGet)
				user.POST("me/apikey", userHandler.CurrentUserApikeyCreate)
			}
		}
	}

	notify(env, absenceRepo)

	r.Run()
}

func migrateExternalCalendarMulti(migrationRepo *repository.Migration, absenceRepo *repository.Absence) error {
	_, err := migrationRepo.MigrationFindByTitle(MIGRATION_EXTERNAL_CALENDAR_MULTI)
	migrationExists := true

	if err != nil {
		if err == repository.ErrMigrationNotFound {
			migrationExists = false
		} else {
			return err
		}
	}

	if migrationExists {
		log.Println("Migration: MIGRATION_EXTERNAL_CALENDAR_MULTI already finished")
		return nil
	}

	log.Println("Migration: MIGRATION_EXTERNAL_CALENDAR_MULTI started")
	absences, err := absenceRepo.FindAll(true)
	if err != nil {
		panic(err)
	}

	for _, absence := range absences {
		if absence.ExternalEventID == "" {
			continue
		}

		eventExists := false
		for _, event := range absence.ExternalEvents {
			if event.ExternalEventID == absence.ExternalEventID {
				eventExists = true
				break
			}
		}

		if eventExists {
			continue
		}

		absenceExternalEvent := model.AbsenceExternalEvent{
			Absence:               absence,
			ExternalEventProvider: absence.ExternalEventProvider,

			ExternalEventID: absence.ExternalEventID,
		}

		err = absenceRepo.AbsenceExternalEventInsert(&absenceExternalEvent)
		if err != nil {
			migration := model.Migration{
				Title:      MIGRATION_EXTERNAL_CALENDAR_MULTI,
				Result:     err.Error(),
				FinishedAt: time.Now(),
				Success:    false,
			}
			migrationRepo.MigrationInsert(&migration)

			return err
		}

		absence.ExternalEventID = "<migrated>"
		absence.ExternalEventProvider = ""

		err = absenceRepo.Update(&absence)
		if err != nil {
			migration := model.Migration{
				Title:      MIGRATION_EXTERNAL_CALENDAR_MULTI,
				Result:     err.Error(),
				FinishedAt: time.Now(),
				Success:    false,
			}
			migrationRepo.MigrationInsert(&migration)

			return err
		}
	}

	migration := model.Migration{
		Title:      MIGRATION_EXTERNAL_CALENDAR_MULTI,
		Result:     "events migrated",
		FinishedAt: time.Now(),
		Success:    true,
	}
	migrationRepo.MigrationInsert(&migration)

	log.Println("Migration: MIGRATION_EXTERNAL_CALENDAR_MULTI finished")
	return nil
}

func migrateExternalCalendar(migrationRepo *repository.Migration, absenceRepo *repository.Absence) error {
	if !microsoft.IsMicrosoftConnected() {
		log.Println("Migration: MIGRATION_EXTERNAL_CALENDAR skipped (no microsoft connection)")
		return nil
	}

	_, err := migrationRepo.MigrationFindByTitle(MIGRATION_EXTERNAL_CALENDAR)
	migrationExists := true

	if err != nil {
		if err == repository.ErrMigrationNotFound {
			migrationExists = false
		} else {
			return err
		}
	}

	if migrationExists {
		log.Println("Migration: MIGRATION_EXTERNAL_CALENDAR already finished")
		return nil
	}

	log.Println("Migration: MIGRATION_EXTERNAL_CALENDAR started")
	absences, err := absenceRepo.FindAll(true)
	if err != nil {
		panic(err)
	}

	for _, absence := range absences {
		if absence.AbsenceFrom.Year() < 2025 {
			continue
		}

		if absence.ExternalEventID != "" {
			continue
		}

		if absence.Identifier == uuid.Nil {
			absence.Identifier = uuid.New()
			absenceRepo.Update(&absence)
		}

		eventId, err := microsoft.CreateCalendarEntryFromAbsence(absence.User.Username, &absence)
		if err != nil {
			return err
		}

		absence.ExternalEventID = eventId
		absence.ExternalEventProvider = model.EXTERNAL_EVENT_PROVIDER_MICROSOFT

		err = absenceRepo.Update(&absence)
		if err != nil {
			migration := model.Migration{
				Title:      MIGRATION_EXTERNAL_CALENDAR,
				Result:     err.Error(),
				FinishedAt: time.Now(),
				Success:    false,
			}
			migrationRepo.MigrationInsert(&migration)

			return err
		}
	}
	migration := model.Migration{
		Title:      MIGRATION_EXTERNAL_CALENDAR,
		Result:     "events migrated",
		FinishedAt: time.Now(),
		Success:    true,
	}
	migrationRepo.MigrationInsert(&migration)

	log.Println("Migration: MIGRATION_EXTERNAL_CALENDAR finished")

	return nil
}

func importHolidays(holiday *repository.Holiday) {
	importHolidaysCurrentYear(holiday)

	for range time.Tick(time.Hour * 24) {
		importHolidaysCurrentYear(holiday)
	}
}

func importHolidaysCurrentYear(holiday *repository.Holiday) {
	year := time.Now().Year()

	log.Printf("Holiday Import: %d", year)
	err := importHolidaysByYear(holiday, year)
	if err != nil {
		log.Println(err)
	}
}

func importHolidaysByYear(holiday *repository.Holiday, year int) error {
	holidays, err := holiday.HolidayFindByYear(year)
	if err != nil {
		panic(err)
	}
	if len(holidays) > 0 {
		return nil
	}

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

		exists, err := holiday.HolidayIsByDate(date)
		if err != nil {
			return err
		}

		if !exists {
			err = holiday.HolidayInsert(&model.Holiday{
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

func notify(env *core.Environment, absenceRepo *repository.Absence) {
	checkIntervalTicker := time.NewTicker(30 * time.Second)
	send := false
	go func() {
		for {
			select {
			case <-checkIntervalTicker.C:
				now := time.Now()
				if now.Weekday() == time.Monday && now.Hour() == 8 && now.Minute() == 0 {
					if !send {
						worker.NotifyAbsenceWeek(env, absenceRepo)
					}
					send = true
				} else {
					send = false
				}
			}
		}
	}()
}
