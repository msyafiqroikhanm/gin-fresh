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

// Custome Error Message For Request / Form Validator.
// If there are error tag missing add new one inside switch case.
func customMessage(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "lg":
		return "Data "
	}
	return "Invalid Field"
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

func ValidationErrorHandlerV1(c *gin.Context, err error) interface{} {
	if errs, ok := err.(validator.ValidationErrors); ok {
		errorMessages := make(map[string]string)
		for _, e := range errs {
			errorMessages[e.Field()] = customMessage(e.Tag())
		}

		return map[string]interface{}{"errors": errorMessages}
	} else {
		return nil
	}
}
