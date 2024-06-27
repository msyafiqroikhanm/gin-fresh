package handlers

import (
	"fmt"
	"jxb-eprocurement/handlers/dtos"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response structure to maintain order
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Standard Service Response To Controller
type ServiceResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Err     interface{} `json:"error"`
}

// Standard Service Response With Logging To Controller
type ServiceResponseWithLogging struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Err     interface{} `json:"error"`
	Log     Log
}

type Log struct {
	Location  string
	StartTime time.Time
	EndTime   time.Time
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

// func ResponseFormatter(c *gin.Context, responseLogging ServiceResponseWithLogging) {
// 	response := Response{
// 		Success: false,
// 		Message: responseLogging.Message,
// 		Data:    responseLogging.Data,
// 	}

// 	if responseLogging.Status == http.StatusOK || responseLogging.Status == http.StatusCreated {
// 		response.Success = true
// 	}
// 	if responseLogging.Data == nil {
// 		// If data is nil, replace it with an empty struct
// 		response.Data = struct{}{}
// 	}

// 	// Check if struct empty
// 	if !responseLogging.Log.StartTime.IsZero() {
// 		fmt.Println(responseLogging.Log)

// 		logSystemParam := LogSystemParam{
// 			Identifier: c.GetString("X-Request-ID"),
// 			StatusCode: responseLogging.Status,
// 			Location:   responseLogging.Log.Location,
// 			Message:    responseLogging.Message,
// 			StartTime:  responseLogging.Log.StartTime,
// 			EndTime:    responseLogging.Log.EndTime,
// 		}

// 		LogSystem(logSystemParam)
// 	}

// 	c.JSON(responseLogging.Status, response)
// }

func ResponseFormatterWithLogging(c *gin.Context, responseLogging ServiceResponseWithLogging) {
	var (
		userLog  = dtos.LogUserInfo{}
		response = Response{
			Success: false,
			Message: responseLogging.Message,
			Data:    responseLogging.Data,
		}
	)

	if responseLogging.Status == http.StatusOK || responseLogging.Status == http.StatusCreated {
		response.Success = true
	}
	if responseLogging.Data == nil {
		// If data is nil, replace it with an empty struct
		response.Data = struct{}{}
	}

	// Get UserID from session
	if sessionUserID, ok := c.Get("UserID"); ok {
		userLog.ID = sessionUserID.(string)
	}

	// Get username from session
	if sessionUsername, ok := c.Get("username"); ok {
		userLog.Username = sessionUsername.(string)
	}

	// Check if struct empty
	if !responseLogging.Log.StartTime.IsZero() {
		fmt.Println(responseLogging.Log)

		logSystemParam := LogSystemParam{
			Identifier: c.GetString("X-Request-ID"),
			StatusCode: responseLogging.Status,
			Location:   responseLogging.Log.Location,
			Message:    responseLogging.Message,
			StartTime:  responseLogging.Log.StartTime,
			EndTime:    responseLogging.Log.EndTime,
			UserInfo:   userLog,
			Err:        responseLogging.Err,
		}

		LogSystem(logSystemParam)
	}

	c.JSON(responseLogging.Status, response)
}
