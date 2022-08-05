package response

import "github.com/gin-gonic/gin"

type Response interface {
	NewErrorResponse(status int, errorCode string, message string, data interface{})
	NewSuccesMessage(status int, message string, data interface{})
}

type response struct {
	ctx *gin.Context
	// StatusErr int
	// Error     bool        `json:"error"`
	// Message   string      `json:"message"`
	// Data      interface{} `json:"data"`
}

func (r *response) NewErrorResponse(status int, errorCode string, message string, data interface{}) {
	r.ctx.JSON(status, gin.H{
		"error":      true,
		"error_code": errorCode,
		"message":    message,
		"data":       nil,
	})
}

func (r *response) NewSuccesMessage(status int, message string, data interface{}) {
	r.ctx.JSON(status, gin.H{
		"error":   false,
		"message": message,
		"data":    data,
	})
}

func NewResponse(c *gin.Context) Response {
	return &response{
		ctx: c,
		// StatusErr: status,
		// Error:     isError,
		// Message:   message,
		// Data:      data,
	}
}
