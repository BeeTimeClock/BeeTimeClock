package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BeeTimeClock/BeeTimeClock-Server/auth"
	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
)

type User struct {
	env  *core.Environment
	user *repository.User
}

func NewUser(env *core.Environment, user *repository.User) *User {
	return &User{
		env:  env,
		user: user,
	}
}

func (h *User) AdministrationUserGetAll(c *gin.Context) {
	users, err := h.user.FindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	var result []model.UserResponse
	for _, user := range users {
		result = append(result, user.GetUserResponse())
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *User) AdministrationUserGetByUserID(c *gin.Context) {
	userIdParam := c.Param("userID")
	userId, err := strconv.ParseUint(userIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("missing userid")))
		return
	}

	user, err := h.user.FindByID(uint(userId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(user.GetUserResponse()))
}

func (h *User) AdministrationUserUpdate(c *gin.Context) {
	var userUpdateRequest model.UserUpdateRequest
	err := c.BindJSON(&userUpdateRequest)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	userIdParam := c.Param("userID")
	userId, err := strconv.ParseUint(userIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("missing userid")))
		return
	}

	user, err := h.user.FindByID(uint(userId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	user.FirstName = userUpdateRequest.FirstName
	user.LastName = userUpdateRequest.LastName
	user.AccessLevel = userUpdateRequest.AccessLevel

	err = h.user.Update(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(user.GetUserResponse()))
}

func (h *User) AdministrationUserCreate(c *gin.Context) {
	var userCreateRequest model.UserCreateRequest
	err := c.BindJSON(&userCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	user := model.NewUser(userCreateRequest.Username)
	user.AccessLevel = userCreateRequest.AccessLevel
	user.FirstName = userCreateRequest.FirstName
	user.LastName = userCreateRequest.LastName

	err = user.SetPassword(userCreateRequest.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	err = h.user.Insert(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(user.GetUserResponse()))
}

func (h *User) AdministrationUserDelete(c *gin.Context) {
	userIDParam := c.Param("userID")

	if strings.TrimSpace(userIDParam) == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("userID missing")))
		return
	}

	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	user, err := h.user.FindByID(uint(userID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err)
		return
	}

	err = h.user.Delete(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *User) CurrentUserGet(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(user.GetUserResponse()))
}
