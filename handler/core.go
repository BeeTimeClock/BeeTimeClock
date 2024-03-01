package handler

import (
	"net/http"
	"strconv"

	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
)

func getUserFromParam(c *gin.Context, userRepo *repository.User) (model.User, bool) {
	userIdParam := c.Param("userID")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return model.User{}, false
	}

	user, err := userRepo.FindByID(uint(userId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return model.User{}, false
	}

	return user, true
}

func getYearFromParam(c *gin.Context) (int, bool) {
	yearParam := c.Param("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return 0, false
	}

	return year, true
}

func getYearMonthFromParam(c *gin.Context) (int, int, bool) {
	year, success := getYearFromParam(c)
	if !success {
		return year, 0, false
	}

	monthParam := c.Param("month")
	month, err := strconv.Atoi(monthParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return year, 0, false
	}

	return year, month, true
}
