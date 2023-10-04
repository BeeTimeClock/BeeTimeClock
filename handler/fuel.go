package handler

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/BeeTimeClock/BeeTimeClock-Server/auth"
	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/helper"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
)

type Fuel struct {
	env  *core.Environment
	user *repository.User
	fuel *repository.Fuel
}

func NewFuel(env *core.Environment, user *repository.User, fuel *repository.Fuel) *Fuel {
	return &Fuel{
		env:  env,
		user: user,
		fuel: fuel,
	}
}

func (h *Fuel) FuelGetAll(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	var fuelQuery model.FuelQueryAll
	err = c.BindQuery(&fuelQuery)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	var fuels []model.Fuel
	if fuelQuery.HasState() {
		fuels, err = h.fuel.FindByUserIDAndState(user.ID, fuelQuery.State)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
			return
		}
	} else {
		fuels, err = h.fuel.FindByUserID(user.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
			return
		}
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(fuels))
}

func (h *Fuel) FuelGet(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	fuelIDParam := c.Param("fuelID")
	fuelID, err := strconv.ParseUint(fuelIDParam, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	fuel, err := h.fuel.FindByID(uint(fuelID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if *fuel.UserID != user.ID {
		c.AbortWithStatusJSON(http.StatusForbidden, model.NewErrorResponse(fmt.Errorf("no rights")))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(fuel))
}

func (h *Fuel) FuelUpdate(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	var fuelUpdateRequest model.FuelUpdateRequest
	err = c.BindJSON(&fuelUpdateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	fuelIDParam := c.Param("fuelID")
	fuelID, err := strconv.ParseUint(fuelIDParam, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	fuel, err := h.fuel.FindByID(uint(fuelID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if *fuel.UserID != user.ID {
		c.AbortWithStatusJSON(http.StatusForbidden, model.NewErrorResponse(fmt.Errorf("no rights")))
		return
	}

	fuel.ReceiptDate = fuelUpdateRequest.ReceiptDate
	fuel.ReceiptValue = fuelUpdateRequest.ReceiptValue
	fuel.ReceiptFuelValue = fuelUpdateRequest.ReceiptFuelValue
	fuel.State = model.FUEL_STATE_OPEN

	err = h.fuel.Update(&fuel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(fuel))
}

func (h *Fuel) FuelActionPrepare(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	fuel := model.Fuel{
		User:  &user,
		State: model.FUEL_STATE_PREPARED,
	}

	err = h.fuel.Insert(&fuel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	uploadPath := path.Join(h.env.UploadPath, "fuel", fmt.Sprintf("%d_%d", user.ID, fuel.ID))
	err = c.SaveUploadedFile(file, uploadPath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}
	fuel.UploadFileName = uploadPath

	text, err := helper.GetTextFromImage(uploadPath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	fuel.ReceiptRawText = text

	err = h.fuel.Update(&fuel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(fuel))
}
