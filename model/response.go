package model

import (
	"fmt"
	"time"
)

type BaseResponse struct {
	Timestamp time.Time
	Status    string
}

type SuccessResponse struct {
	BaseResponse
	Data interface{}
}

type ErrorResponse struct {
	BaseResponse
	Error   error
	Message string
}

func getBaseResponse(status string) BaseResponse {
	return BaseResponse{
		Timestamp: time.Now(),
		Status:    status,
	}
}

func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{
		BaseResponse: getBaseResponse("success"),
		Data:         data,
	}
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		BaseResponse: getBaseResponse("error"),
		Error:        err,
		Message:      err.Error(),
	}
}

func NewMessageResponse(format string, a ...any) SuccessResponse {
	type Message struct {
		Message string
	}

	return NewSuccessResponse(Message{
		Message: fmt.Sprintf(format, a...),
	})
}

type SumResult struct {
	Total float64
}
