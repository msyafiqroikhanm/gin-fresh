package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(data interface{}) error {
	return validate.Struct(data)
}

func ValidationErrorHandler(c *gin.Context, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, e := range errs {
			errorMessages = append(errorMessages, e.Error())
		}
		ResponseFormatter(c, http.StatusBadRequest, nil, errorMessages[0])
	} else {
		ResponseFormatter(c, http.StatusInternalServerError, nil, err.Error())
	}
}
