package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jxb-eprocurement/handlers/dtos"
	"net/url"
	"os"
	"strconv"
	"strings"
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

	// fmt.Println("Init Logger")

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

func APILogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			startTime  = time.Now()
			identifier = c.GetString("X-Request-ID")
		)

		// Read the request body
		var requestBody interface{}
		contentType := c.Request.Header.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			err := c.Request.ParseMultipartForm(32 << 20) // 32 MB
			if err != nil {
				requestBody = fmt.Sprintf("Error parsing form-data: %v", err)
			} else {
				multipartData := make(map[string]interface{})
				for key, values := range c.Request.MultipartForm.Value {
					if key == "password" || key == "re_password" || key == "old_password" {
						multipartData[key] = "[REDACTED]"
					} else {
						if len(values) > 1 {
							multipartData[key] = values
						} else {
							multipartData[key] = values[0]
						}
					}
				}
				requestBody = multipartData
			}
		} else {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			if strings.Contains(contentType, "application/x-www-form-urlencoded") {
				formData, err := url.ParseQuery(string(bodyBytes))
				if err != nil {
					requestBody = string(bodyBytes)
				} else {
					jsonFormData := make(map[string]interface{})
					for key, values := range formData {
						if key == "password" || key == "re_password" || key == "old_password" {
							jsonFormData[key] = "[REDACTED]"
						} else {
							if len(values) > 1 {
								jsonFormData[key] = values
							} else {
								jsonFormData[key] = values[0]
							}
						}
					}
					requestBody = jsonFormData
				}
			} else {
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

		// Redact Authorization header
		headers := c.Request.Header.Clone()
		if _, ok := headers["Authorization"]; ok {
			headers["Authorization"] = []string{"[REDACTED]"}
		}

		var (
			endTime      = time.Now()
			userAgent    = c.Request.UserAgent()
			clientIP     = c.ClientIP()
			endpoint     = c.Request.URL.Path
			queryParams  = c.Request.URL.Query()
			statusCode   = c.Writer.Status()
			responseBody = responseWriter.body.Bytes()
			httpMethod   = c.Request.Method
		)

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

		// Get username from session
		var username string
		if sessionUsername, ok := c.Get("username"); ok {
			username = sessionUsername.(string)
		}

		var logFunc func(string, ...zapcore.Field)
		switch {
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
			zap.String("http_method", httpMethod),
			zap.Any("request_header", headers),
			zap.Any("query_params", queryParams),
			zap.Any("request_body", requestBody),
			zap.String("response_code", statusCodeString),
			zap.Any("response_body", jsonResponseBody),
			zap.String("endpoint", endpoint),
			zap.String("user_agent", userAgent),
			zap.String("client_ip", clientIP),
			zap.String("username", username),
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
	UserInfo   dtos.LogUserInfo
	Err        interface{}
}

func LogSystem(logData LogSystemParam) {
	var (
		userInfo         interface{}
		category         string
		humanTime        = logData.EndTime.Format(time.RFC1123)
		statusCodeString = strconv.Itoa(logData.StatusCode)
	)

	switch true {
	case logData.StatusCode == 500:
		category = "FATAL"
	case logData.StatusCode < 500 && logData.StatusCode >= 400:
		category = "ERROR"
	default:
		category = "INFO"
	}

	if logData.UserInfo.ID == "" && logData.UserInfo.Username == "" {
		userInfo = struct{}{}
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
		zap.Any("user_info", userInfo),
		zap.Any("errors", logData.Err),
		zap.String("human_time", humanTime),
	)
}
