package mongolia

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    int    `json:"code" bson:"code"`
	Message string `json:"message" bson:"message"`
	Error   error  `json:"error" bson:"error"`
}

func (e *Error) Set(code int, message string) {
	e.Code = code
	e.Message = message
	e.Error = errors.New(message)
}

func (e *Error) ToString() string {
	return fmt.Sprintf("Code: %v Message: %v ", e.Code, e.Message)
}

func NewError(code int, err error) *Error {
	e := Error{
		Code:    code,
		Message: err.Error(),
		Error:   err,
	}
	return &e
}

func NewErrorString(code int, message string) *Error {
	e := Error{
		Code:    code,
		Message: message,
		Error:   errors.New(message),
	}
	return &e
}
