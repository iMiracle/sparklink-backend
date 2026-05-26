package response

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	ErrorCode string      `json:"errorCode,omitempty"`
	RequestId string      `json:"requestId"`
}

func generateRequestId() string {
	b := make([]byte, 8)
	rand.Read(b)
	return "req_" + hex.EncodeToString(b)
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ApiResponse{
		Code:      0,
		Message:   "success",
		Data:      data,
		RequestId: generateRequestId(),
	})
}

func Error(c *gin.Context, httpStatus int, errCode string, message string) {
	if message == "" {
		if msg, ok := ErrorMessages[errCode]; ok {
			message = msg
		}
	}
	c.JSON(httpStatus, ApiResponse{
		Code:      httpStatus,
		Message:   message,
		ErrorCode: errCode,
		RequestId: generateRequestId(),
	})
}

func BadRequest(c *gin.Context, errCode string, message string) {
	Error(c, http.StatusBadRequest, errCode, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, ErrAuthFailed, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, ErrPermissionDenied, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, ErrNotFound, message)
}

func TooManyRequests(c *gin.Context, message string) {
	Error(c, http.StatusTooManyRequests, ErrRateLimited, message)
}

func ServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, ErrServerError, message)
}
