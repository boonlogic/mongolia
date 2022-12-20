package mongolia

import (
	"fmt"
)

type Error struct {
	Code    int    `json:"code" bson:"code"`
	Message string `json:"message" bson:"message"`
}

func (e *Error) Set(code int, message string) {
	e.Code = code
	e.Message = message
}

func (e *Error) ToString() string {
	return fmt.Sprintf("Code: %v Message: %v ", e.Code, e.Message)
}

func NewError(code int, err error) *Error {
	e := Error{
		Code:    code,
		Message: err.Error(),
	}
	return &e
}

func NewErrorString(code int, message string) *Error {
	e := Error{
		Code:    code,
		Message: message,
	}
	return &e
}
