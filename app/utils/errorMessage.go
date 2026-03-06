package utils

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

// Generic FieldError handling. For custom handling of specific fields or structs, define your own implementation.
type ValidationFieldError struct {
	Err validator.FieldError
}

// String generates an error message based on the validation error's tag.
// Supported tags include "required", "max", "min", "email", "len", "gt", "gte", "lt", "lte", "oneof".
// For unknown tags, a default error message format is returned.
// The returned error message includes the field name and its corresponding condition.
func (v ValidationFieldError) String() string {
	e := v.Err

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", e.Field(), e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", e.Field(), e.Param())
	case "email":
		return "Invalid email format"
	case "len":
		return fmt.Sprintf("%s must be %s characters long", e.Field(), e.Param())
	case "gt":
		return fmt.Sprintf("%s must greater than %s", e.Field(), e.Param())
	case "gte":
		return fmt.Sprintf("%s must greater or equals to %s", e.Field(), e.Param())
	case "lt":
		return fmt.Sprintf("%s must less than %s", e.Field(), e.Param())
	case "lte":
		return fmt.Sprintf("%s must less or equals to %s", e.Field(), e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of '%s'", e.Field(), e.Param())
	}

	return fmt.Sprintf("%s is not valid, condition: %s", e.Field(), e.ActualTag())
}

// ValidationErrorMessage returns the corresponding validation error message based on the provided error.
// If the error is io.EOF, returns "EOF, json decode fail".
// If the error is validator.ValidationErrors, returns the message of the first validation error.
// If the error is not validator.ValidationErrors, returns "json decode or validate fail, err=" plus the error message.
// If there is no error message, returns "validationErrs with no error message".
func ValidationErrorMessage(err error) string {
	if err == io.EOF {
		return "EOF, json decode fail"
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		message := fmt.Sprintf("json decode or validate fail, err=%s", err)
		log.Info(message)
		return message
	}

	// currently, only return the first error
	for _, fieldErr := range validationErrs {
		return ValidationFieldError{fieldErr}.String()
	}

	return "validationErrs with no error message"
}

// HandleError is a generic error handling function
func HandleError(c *gin.Context, status int, code int, msg string, err error) {
	if err != nil {
		msg = msg + ": " + err.Error()
	}

	c.JSON(status, gin.H{
		"code": code,
		"msg":  msg,
		"data": nil,
	})
}
