package tools

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// CodeSuccess ...
	CodeSuccess = 0
	// CodeFail ...
	CodeFail = 1
	// CodeUnknownError ...
	CodeUnknownError = -1
	// CodeSessionError ...
	CodeSessionError = 40000
)

// MsgCodeMap ...
var MsgCodeMap = map[int]string{
	CodeUnknownError: "unKnow error",
	CodeSuccess:      "success",
	CodeFail:         "fail",
	CodeSessionError: "Session error",
}

// SuccessWithMsg ...
func SuccessWithMsg(c *gin.Context, msg interface{}, data interface{}) {
	ResponseWithCode(c, CodeSuccess, msg, data)
}

// FailWithMsg ...
func FailWithMsg(c *gin.Context, msg interface{}) {
	ResponseWithCode(c, CodeFail, msg, nil)
}

// ResponseWithCode ...
func ResponseWithCode(c *gin.Context, msgCode int, msg interface{}, data interface{}) {
	if msg == nil {
		if val, ok := MsgCodeMap[msgCode]; ok {
			msg = val
		} else {
			msg = MsgCodeMap[CodeUnknownError]
		}
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":    msgCode,
		"message": msg,
		"data":    data,
	})
}
