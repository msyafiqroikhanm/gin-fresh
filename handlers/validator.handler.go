package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("no_space", noSpace)
}

func ValidateStruct(data interface{}) error {
	return validate.Struct(data)
}

// Custome Error Message For Request / Form Validator.
// If there are error tag missing add new one inside switch case.
func customMessage(tag string, size string) string {
	switch tag {
	case "required":
		return "Field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Field require minimum of %s size/length/unit", size)
	case "no_space":
		return "Field should not contain spaces"
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

func ValidationErrorHandlerV1(c *gin.Context, err error, dto interface{}) interface{} {
	if errs, ok := err.(validator.ValidationErrors); ok {
		errorMessages := make(map[string]string)

		// Use reflection to map struct field names to JSON tag names
		dtoType := reflect.TypeOf(dto)
		for _, e := range errs {
			// Get the struct field
			field, _ := dtoType.FieldByName(e.Field())
			// Get the JSON tag
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" {
				jsonTag = e.Field()
			}
			fmt.Println(e.Tag())
			errorMessages[jsonTag] = customMessage(e.Tag(), e.Param())
		}

		return map[string]interface{}{"errors": errorMessages}
	} else {
		return nil
	}
}

// Custom validator that validate given input doens't have any spaces
func noSpace(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return !strings.Contains(value, " ")
}
