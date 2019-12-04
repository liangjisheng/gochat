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

// Login ...
func Login(c *gin.Context) {
	var formLogin FormLogin
	if err := c.ShouldBindBodyWith(&formLogin, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	req := &proto.LoginRequest{
		Name:     formLogin.UserName,
		Password: tools.Sha1(formLogin.Password),
	}
	code, authToken, msg := rpc.RPCLogicObj.Login(req)
	if code == tools.CodeFail || authToken == "" {
		tools.FailWithMsg(c, msg)
		return
	}
	tools.SuccessWithMsg(c, "login success", authToken)
}

// FormCheckAuth ...
type FormCheckAuth struct {
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

// CheckAuth ...
func CheckAuth(c *gin.Context) {
	var formCheckAuth FormCheckAuth
	if err := c.ShouldBindBodyWith(&formCheckAuth, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	authToken := formCheckAuth.AuthToken
	req := &proto.CheckAuthRequest{
		AuthToken: authToken,
	}
	code, userID, userName := rpc.RPCLogicObj.CheckAuth(req)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, "auth fail")
		return
	}
	var jsonData = map[string]interface{}{
		"userId":   userID,
		"userName": userName,
	}
	tools.SuccessWithMsg(c, "auth success", jsonData)
}

// FormLogout ...
type FormLogout struct {
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

// Logout ...
func Logout(c *gin.Context) {
	var formLogout FormLogout
	if err := c.ShouldBindBodyWith(&formLogout, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	authToken := formLogout.AuthToken
	logoutReq := &proto.LogoutRequest{
		AuthToken: authToken,
	}
	code := rpc.RPCLogicObj.Logout(logoutReq)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, "logout fail!")
		return
	}
	tools.SuccessWithMsg(c, "logout ok!", nil)
}
