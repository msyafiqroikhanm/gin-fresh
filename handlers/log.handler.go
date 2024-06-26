package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	apiLogger    *zap.Logger
	systemLogger *zap.Logger
)

func InitLogger() {
	// Create logs directory if it does not exist
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		panic(err)
	}

	fmt.Println("Init Logger")

	// Encoder configuration
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Console output configuration
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// API Logger setup
	apiLogFile, err := os.OpenFile("logs/api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	apiFileWriter := zapcore.AddSync(apiLogFile)
	apiCore := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), apiFileWriter, zapcore.InfoLevel),
	)
	apiLogger = zap.New(apiCore)

	// System Logger setup
	systemLogFile, err := os.OpenFile("logs/system.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	systemFileWriter := zapcore.AddSync(systemLogFile)
	systemCore := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), systemFileWriter, zapcore.InfoLevel),
	)
	systemLogger = zap.New(systemCore)
}

// func RequestIDMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		requestID := c.Request.Header.Get("X-Request-ID")
// 		if requestID == "" {
// 			requestID = uuid.New().String()
// 			c.Request.Header.Set("X-Request-ID", requestID)
// 		}
// 		c.Writer.Header().Set("X-Request-ID", requestID)
// 		c.Set("X-Request-ID", requestID)
// 		c.Next()
// 	}
// }

func APILogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		identifier := c.GetString("X-Request-ID")

		/// Read the request body
		var requestBody interface{}
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)

			// Parse body based on content type
			contentType := c.Request.Header.Get("Content-Type")
			if contentType == "application/x-www-form-urlencoded" {
				// Convert URL encoded body to JSON
				formData, err := url.ParseQuery(string(bodyBytes))
				if err != nil {
					requestBody = string(bodyBytes)
				} else {
					requestBody = formData
				}
			} else {
				// Assume JSON or other format
				err := json.Unmarshal(bodyBytes, &requestBody)
				if err != nil {
					requestBody = string(bodyBytes)
				}
			}

			// Restore the request body for downstream handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Create a custom response writer
		responseWriter := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = responseWriter

		c.Next()

		endTime := time.Now()
		userAgent := c.Request.UserAgent()
		clientIP := c.ClientIP()
		endpoint := c.Request.URL.Path
		headers := c.Request.Header
		queryParams := c.Request.URL.Query()

		responseBody := responseWriter.body.Bytes()

		var jsonResponseBody map[string]interface{}
		if err := json.Unmarshal(responseBody, &jsonResponseBody); err != nil {
			// If the response is not JSON, log it as a string
			jsonResponseBody = map[string]interface{}{
				"raw": string(responseBody),
			}
		}

		// Extract message from JSON response body
		var message string
		if msg, ok := jsonResponseBody["message"]; ok {
			message = fmt.Sprintf("API Log | %s", msg.(string))
		} else {
			message = "API LOG | No message in response"
		}

		statusCode := c.Writer.Status()

		// // Get UserID from session
		// var userID string
		// if sessionUserID, ok := c.Get("UserID"); ok {
		// 	userID = sessionUserID.(string)
		// }

		// Get username from session
		var username string
		if sessionUsername, ok := c.Get("username"); ok {
			username = sessionUsername.(string)
		}

		var logFunc func(string, ...zapcore.Field)
		switch true {
		case statusCode == 500:
			logFunc = apiLogger.Fatal
		case statusCode < 500 && statusCode >= 400:
			logFunc = apiLogger.Error
		default:
			logFunc = apiLogger.Info
		}

		statusCodeString := strconv.Itoa(statusCode)

		logFunc(message,
			zap.String("identifier", identifier),
			zap.Time("timestamp", time.Now()),
			zap.Any("request_header", headers),
			zap.Any("query_params", queryParams),
			zap.Any("request_body", requestBody),
			zap.String("response_code", statusCodeString),
			zap.Any("response_body", jsonResponseBody),
			zap.String("endpoint", endpoint),
			zap.String("user_agent", userAgent),
			zap.String("client_ip", clientIP),
			zap.String("username", username),
			// zap.String("user_id", userID),
			zap.Time("start_time", startTime),
			zap.Time("end_time", endTime),
		)
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

type LogSystemParam struct {
	Identifier string
	StatusCode int
	Location   string
	Message    string
	StartTime  time.Time
	EndTime    time.Time
}

func LogSystem(logData LogSystemParam) {
	humanTime := logData.EndTime.Format(time.RFC1123)
	statusCodeString := strconv.Itoa(logData.StatusCode)

	var category string
	switch true {
	case logData.StatusCode == 500:
		category = "FATAL"
	case logData.StatusCode < 500 && logData.StatusCode >= 400:
		category = "ERROR"
	default:
		category = "INFO"
	}

	systemLogger.Info("System Log",
		zap.Time("timestamp", time.Now()),
		zap.String("category", category),
		zap.String("response_code", statusCodeString),
		zap.String("location", logData.Location),
		zap.String("message", logData.Message),
		zap.Time("start_time", logData.StartTime),
		zap.Time("end_time", logData.EndTime),
		zap.String("identifier", logData.Identifier),
		zap.String("human_time", humanTime),
	)
}
