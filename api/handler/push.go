package handler

import (
	"gochat/api/rpc"
	"gochat/config"
	"gochat/proto"
	"gochat/tools"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// FormPush ...
type FormPush struct {
	Msg       string `form:"msg" json:"msg" binding:"required"`
	ToUserID  string `form:"toUserId" json:"toUserId" binding:"required"`
	RoomID    int    `form:"roomId" json:"roomId" binding:"required"`
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

// Push ...
func Push(c *gin.Context) {
	var formPush FormPush
	if err := c.ShouldBindBodyWith(&formPush, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	authToken := formPush.AuthToken
	msg := formPush.Msg
	toUserID := formPush.ToUserID
	toUserIDInt, _ := strconv.Atoi(toUserID)
	getUserNameReq := &proto.GetUserInfoRequest{UserID: toUserIDInt}
	code, toUserName := rpc.RPCLogicObj.GetUserNameByUserID(getUserNameReq)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, "rpc fail get friend userName")
		return
	}
	checkAuthReq := &proto.CheckAuthRequest{AuthToken: authToken}
	code, fromUserID, fromUserName := rpc.RPCLogicObj.CheckAuth(checkAuthReq)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, "rpc fail get self info")
		return
	}
	roomID := formPush.RoomID
	req := &proto.Send{
		Msg:          msg,
		FromUserID:   fromUserID,
		FromUserName: fromUserName,
		ToUserID:     toUserIDInt,
		ToUserName:   toUserName,
		RoomID:       roomID,
		Op:           config.OpSingleSend,
	}
	code, rpcMsg := rpc.RPCLogicObj.Push(req)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, rpcMsg)
		return
	}
	tools.SuccessWithMsg(c, "ok", nil)
	return
}
