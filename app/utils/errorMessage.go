package utils

import (
	"fmt"
	"io"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

// 這裡是通用的 FieldError 處理, 如果需要針對某些字段或 struct 做定制, 需要自行定義一個
type ValidationFieldError struct {
	Err validator.FieldError
}


// String 會根據驗證錯誤的標籤 (Tag) 生成對應的錯誤訊息。
// 支援的標籤包括 "required", "max", "min", "email", "len", "gt", "gte", "lt", "lte", "oneof"。
// 對於未知的標籤，會返回預設的錯誤訊息格式。
// 返回的錯誤訊息會包含欄位名稱及其對應的條件。
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

// ValidationErrorMessage 根據提供的錯誤訊息返回對應的驗證錯誤訊息。
// 如果錯誤是 io.EOF，返回 "EOF, json decode fail"。
// 如果錯誤是 validator.ValidationErrors，返回第一個驗證錯誤的訊息。
// 如果錯誤不是 validator.ValidationErrors，返回 "json decode or validate fail, err=" 加上錯誤訊息。
// 如果沒有錯誤訊息，返回 "validationErrs with no error message"。
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