package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response structure to maintain order
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ServiceResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Err     error       `json:"error"`
}

func ResponseFormatter(c *gin.Context, status int, data interface{}, message string) {
	response := Response{
		Success: false,
		Message: message,
		Data:    data,
	}
	if status == http.StatusOK || status == http.StatusCreated {
		response.Success = true
	}
	if data == nil {
		// If data is nil, replace it with an empty struct
		response.Data = struct{}{}
	}

	c.JSON(status, response)
}
