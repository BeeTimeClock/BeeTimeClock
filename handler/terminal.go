package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BeeTimeClock/BeeTimeClock-Server/auth"
	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/helper"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
)

type Terminal struct {
	env       *core.Environment
	user      *repository.User
	terminal  *repository.Terminal
	timestamp *repository.Timestamp
}

func NewTerminal(env *core.Environment, user *repository.User, terminal *repository.Terminal, timestamp *repository.Timestamp) *Terminal {
	return &Terminal{
		env:       env,
		user:      user,
		terminal:  terminal,
		timestamp: timestamp,
	}
}

func (h *Terminal) AdministrationTerminalList(c *gin.Context) {
	terminals, err := h.terminal.TerminalFindAll()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(terminals))
}

func (h *Terminal) AdministrationTerminalGet(c *gin.Context) {
	terminalId, err := getIdFromParam(c, "terminalId")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	terminals, err := h.terminal.TerminalFindById(terminalId)

	if err != nil {
		if errors.Is(err, repository.ErrTerminalNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, model.NewErrorResponse(err))
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(terminals))
}

func (h *Terminal) AdministrationTerminalCreate(c *gin.Context) {
	var terminalCreateRequest model.TerminalCreateRequest

	err := c.BindJSON(&terminalCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	terminal := model.Terminal{
		TerminalName: terminalCreateRequest.TerminalName,
		ClientId:     helper.RandomString(32),
	}

	key, err := terminal.GenerateApikey()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	err = h.terminal.TerminalInsert(&terminal)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	terminal.Apikey = key

	c.JSON(http.StatusCreated, model.NewSuccessResponse(terminal))
}

func (h *Terminal) AdministrationTerminalRegenerate(c *gin.Context) {
	terminalId, err := getIdFromParam(c, "terminalId")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	terminal, err := h.terminal.TerminalFindById(terminalId)

	if err != nil {
		if errors.Is(err, repository.ErrTerminalNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, model.NewErrorResponse(err))
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	key, err := terminal.GenerateApikey()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	err = h.terminal.TerminalUpdate(&terminal)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	terminal.Apikey = key

	c.JSON(http.StatusOK, model.NewSuccessResponse(terminal))
}

func (h *Terminal) AdministrationTerminalDelete(c *gin.Context) {
	terminalId, err := getIdFromParam(c, "terminalId")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	terminal, err := h.terminal.TerminalFindById(terminalId)

	if err != nil {
		if errors.Is(err, repository.ErrTerminalNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, model.NewErrorResponse(err))
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	err = h.terminal.TerminalDelete(&terminal)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Terminal) AdministrationUserList(c *gin.Context) {
	userID, err := getIdFromParam(c, "userID")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	tokens, err := h.user.UserTokenFindAllByUserID(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(tokens))
}

func (h *Terminal) AdministrationUserCreate(c *gin.Context) {
	userID, err := getIdFromParam(c, "userID")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	user, err := h.user.FindByID(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, model.NewErrorResponse(err))
		return
	}

	var tokenCreateRequest model.UserTokenCreateRequest
	err = c.BindJSON(&tokenCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	tokenType := tokenCreateRequest.TokenType
	if tokenType == "" {
		tokenType = model.USER_TOKEN_TYPE_CHIP
	}

	userToken := model.UserToken{
		UserID:          user.ID,
		TokenType:       tokenType,
		TokenIdentifier: tokenCreateRequest.TokenIdentifier,
	}

	err = h.user.UserTokenInsert(&userToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(userToken))
}

func (h *Terminal) AdministrationUserTokenDelete(c *gin.Context) {
	userID, err := getIdFromParam(c, "userID")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	tokenId, err := getIdFromParam(c, "tokenId")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	userToken, err := h.user.UserTokenFindById(tokenId)
	if err != nil {
		if errors.Is(err, repository.ErrUserTokenNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, model.NewErrorResponse(err))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if userToken.UserID != userID {
		c.AbortWithStatusJSON(http.StatusNotFound, model.NewErrorResponse(repository.ErrUserTokenNotFound))
		return
	}

	err = h.user.UserTokenDelete(&userToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Terminal) TerminalCheckin(c *gin.Context) {
	var checkinRequest model.TerminalCheckinRequest

	err := c.BindJSON(&checkinRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	device, err := auth.GetDeviceFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	log.Printf("Looking for: %s", checkinRequest.TokenIdentifier)
	userToken, err := h.user.UserTokenFindByTokenIdentifier(checkinRequest.TokenIdentifier)
	if err != nil {
		if errors.Is(err, repository.ErrUserTokenNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, model.NewErrorResponse(err))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	timestamp := model.Timestamp{
		UserID:          userToken.UserID,
		ComingTimestamp: time.Now(),
		CominigDevice:   &device,
		IsHomeoffice:    false,
	}

	err = h.timestamp.Insert(&timestamp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	timestampResponse := model.TerminalTimestampResponse{
		Timestamp: timestamp,
		User:      userToken.User,
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(timestampResponse))
}

func (h *Terminal) TerminalCheckout(c *gin.Context) {
	var checkoutRequest model.TerminalCheckoutRequest

	err := c.BindJSON(&checkoutRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	device, err := auth.GetDeviceFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	userToken, err := h.user.UserTokenFindByTokenIdentifier(checkoutRequest.TokenIdentifier)
	if err != nil {
		if errors.Is(err, repository.ErrUserTokenNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, model.NewErrorResponse(err))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	lastTimestamp, err := h.timestamp.FindLastByUserID(userToken.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if lastTimestamp.IsComplete() {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("there is no open timestamp")))
		return
	}

	lastTimestamp.GoingTimestamp = time.Now()
	lastTimestamp.GoingDevice = &device
	lastTimestamp.IsHomeofficeGoing = false

	err = h.timestamp.Update(&lastTimestamp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	timestampResponse := model.TerminalTimestampResponse{
		Timestamp: lastTimestamp,
		User:      userToken.User,
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(timestampResponse))
}
