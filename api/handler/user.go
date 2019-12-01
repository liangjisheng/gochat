package handler

import (
	"gochat/api/rpc"
	"gochat/proto"
	"gochat/tools"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// FormRegister ...
type FormRegister struct {
	UserName string `form:"userName" json:"userName" binding:"required"`
	Password string `form:"passWord" json:"passWord" binding:"required"`
}

// Register ...
func Register(c *gin.Context) {
	var formRegister FormRegister
	if err := c.ShouldBindBodyWith(&formRegister, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	req := &proto.RegisterRequest{
		Name:     formRegister.UserName,
		Password: tools.Sha1(formRegister.Password),
	}
	code, authToken, msg := rpc.RPCLogicObj.Register(req)
	if code == tools.CodeFail || authToken == "" {
		tools.FailWithMsg(c, msg)
		return
	}
	tools.SuccessWithMsg(c, "register success", authToken)
}

// FormLogin ...
type FormLogin struct {
	UserName string `form:"userName" json:"userName" binding:"required"`
	Password string `form:"passWord" json:"passWord" binding:"required"`
}
