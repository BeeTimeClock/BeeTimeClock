package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/helper"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/BeeTimeClock/BeeTimeClock-Server/worker"
	"github.com/gin-gonic/gin"
)

type Administration struct {
	env      *core.Environment
	settings *repository.Settings
	absence  *repository.Absence
	holiday  *repository.Holiday
}

func NewAdministration(env *core.Environment, settings *repository.Settings, absence *repository.Absence, holiday *repository.Holiday) *Administration {
	return &Administration{
		env:      env,
		settings: settings,
		absence:  absence,
		holiday:  holiday,
	}
}

func (h Administration) AdministrationGetSettings(c *gin.Context) {
	settings, err := h.settings.SettingsFind()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(settings))
}

func (h Administration) AdministrationUpdateSettings(c *gin.Context) {
	var settings model.Settings

	err := c.BindJSON(&settings)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	h.settings.SettingsUpdate(&settings)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(settings))
}

func (h Administration) AdministrationNotifyAbsenceWeek(c *gin.Context) {
	worker.NotifyAbsenceWeek(h.env, h.absence)
	c.Status(http.StatusNoContent)
}

func (h Administration) AdministrationUploadLogo(c *gin.Context) {
	file, _ := c.FormFile("file")
	currentFile, err := file.Open()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	err = helper.SaveFile(h.env, "logo.png", currentFile, file)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h Administration) GetLogo(c *gin.Context) {
	logoName := "logo.png"
	exists, err := helper.ExistsFile(h.env, logoName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if !exists {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	file, stat, err := helper.GetFile(h.env, logoName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`inline; filename="%s"`, logoName),
	}

	c.DataFromReader(http.StatusOK, stat.Size, "application/pdf", file, extraHeaders)
}

func (h Administration) AdministrationGetHolidaysYear(c *gin.Context) {
	year, err := strconv.ParseInt(c.Param("year"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	days, err := h.holiday.HolidayFindByYear(int(year))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(days))
}

func (h Administration) AdministrationGetHolidaysCustom(c *gin.Context) {
	customs, err := h.holiday.HolidayCustomFindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(customs))
}

func (h Administration) AdministrationCreateHolidaysCustom(c *gin.Context) {
	var customHolidayCreateRequest model.HolidayCustom

	err := c.BindJSON(&customHolidayCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	customHoliday := model.HolidayCustom{
		Name:   customHolidayCreateRequest.Name,
		Date:   customHolidayCreateRequest.Date,
		Month:  customHolidayCreateRequest.Month,
		Day:    customHolidayCreateRequest.Day,
		Yearly: customHolidayCreateRequest.Yearly,
	}

	err = h.holiday.HolidayCustomInsert(&customHoliday)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(customHoliday))
}

func (h Administration) AdministrationDeleteHolidaysCustom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	customHoliday, err := h.holiday.HolidayCustomFindById(uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	err = h.holiday.HolidayCustomDelete(&customHoliday)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}
