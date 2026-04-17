package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

type Response struct {
	Code    int          `json:"code"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data,omitempty"`
	Errors  []FieldError `json:"errors,omitempty"`
}

func send(c *gin.Context, code int, success bool, message string, data interface{}, errs []FieldError) {
	c.JSON(code, Response{
		Code:    code,
		Success: success,
		Message: message,
		Data:    data,
		Errors:  errs,
	})
}

func OK(c *gin.Context, message string, data interface{}) {
	send(c, http.StatusOK, true, message, data, nil)
}

func Created(c *gin.Context, message string, data interface{}) {
	send(c, http.StatusCreated, true, message, data, nil)
}

func BadRequest(c *gin.Context, message string) {
	send(c, http.StatusBadRequest, false, message, nil, nil)
}

func Unauthorized(c *gin.Context, message string) {
	send(c, http.StatusUnauthorized, false, message, nil, nil)
}

func InternalServerError(c *gin.Context, message string) {
	send(c, http.StatusInternalServerError, false, message, nil, nil)
}

func Forbidden(c *gin.Context, message string) {
	send(c, http.StatusForbidden, false, message, nil, nil)
}

func NotFound(c *gin.Context, message string) {
	send(c, http.StatusNotFound, false, message, nil, nil)
}

func Conflict(c *gin.Context, message string) {
	send(c, http.StatusConflict, false, message, nil, nil)
}

func ValidationError(c *gin.Context, err error) {
	var details []FieldError

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			details = append(details, FieldError{
				Field:   strings.ToLower(fe.Field()),
				Tag:     fe.Tag(),
				Message: getErrorMsg(fe),
			})
		}
	}

	send(c, http.StatusBadRequest, false, "Validasi gagal", nil, details)
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s wajib diisi", strings.ToLower(fe.Field()))
	case "email":
		return "format email tidak valid"
	case "min":
		return fmt.Sprintf("%s minimal %s karakter", strings.ToLower(fe.Field()), fe.Param())
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter", strings.ToLower(fe.Field()), fe.Param())
	case "alphanum":
		return fmt.Sprintf("%s hanya boleh berisi huruf dan angka", strings.ToLower(fe.Field()))
	}
	return fmt.Sprintf("%s tidak valid", strings.ToLower(fe.Field()))
}
