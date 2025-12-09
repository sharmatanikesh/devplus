package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIError represents a standard API error response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// SendError sends a JSON error response
func SendError(c *gin.Context, code int, message string, err error) {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	c.AbortWithStatusJSON(code, APIError{
		Code:    code,
		Message: message,
		Error:   errStr,
	})
}

// Common Errors
func SendBadRequest(c *gin.Context, err error) {
	SendError(c, http.StatusBadRequest, "Bad Request", err)
}

func SendInternalServerError(c *gin.Context, err error) {
	SendError(c, http.StatusInternalServerError, "Internal Server Error", err)
}

func SendNotFound(c *gin.Context, message string) {
	SendError(c, http.StatusNotFound, message, nil)
}
