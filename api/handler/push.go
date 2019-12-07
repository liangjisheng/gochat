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

// FormRoom ...
type FormRoom struct {
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
	Msg       string `form:"msg" json:"msg" binding:"required"`
	RoomID    int    `form:"roomId" json:"roomId" binding:"required"`
}

// PushRoom ...
func PushRoom(c *gin.Context) {
	var formRoom FormRoom
	if err := c.ShouldBindBodyWith(&formRoom, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	authToken := formRoom.AuthToken
	msg := formRoom.Msg
	roomID := formRoom.RoomID
	checkAuthReq := &proto.CheckAuthRequest{AuthToken: authToken}
	authCode, fromUserID, fromUserName := rpc.RPCLogicObj.CheckAuth(checkAuthReq)
	if authCode == tools.CodeFail {
		tools.FailWithMsg(c, "rpc fail get self info")
		return
	}
	req := &proto.Send{
		Msg:          msg,
		FromUserID:   fromUserID,
		FromUserName: fromUserName,
		RoomID:       roomID,
		Op:           config.OpRoomSend,
	}
	code, msg := rpc.RPCLogicObj.PushRoom(req)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, "rpc push room msg fail!")
		return
	}
	tools.SuccessWithMsg(c, "ok", msg)
	return
}

// FormCount ...
type FormCount struct {
	RoomID int `form:"roomId" json:"roomId" binding:"required"`
}

// Count ...
func Count(c *gin.Context) {
	var formCount FormCount
	if err := c.ShouldBindBodyWith(&formCount, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	roomID := formCount.RoomID
	req := &proto.Send{
		RoomID: roomID,
		Op:     config.OpRoomCountSend,
	}
	code, msg := rpc.RPCLogicObj.Count(req)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, "rpc get room count fail!")
		return
	}
	tools.SuccessWithMsg(c, "ok", msg)
	return
}

// FormRoomInfo ...
type FormRoomInfo struct {
	RoomID int `form:"roomId" json:"roomId" binding:"required"`
}

// GetRoomInfo ...
func GetRoomInfo(c *gin.Context) {
	var formRoomInfo FormRoomInfo
	if err := c.ShouldBindBodyWith(&formRoomInfo, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	roomID := formRoomInfo.RoomID
	req := &proto.Send{
		RoomID: roomID,
		Op:     config.OpRoomInfoSend,
	}
	code, msg := rpc.RPCLogicObj.GetRoomInfo(req)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, "rpc get room info fail!")
		return
	}
	tools.SuccessWithMsg(c, "ok", msg)
	return
}
