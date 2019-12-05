package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gochat/tools"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"gochat/config"
	"gochat/logic/dao"
	"gochat/proto"
)

// RPCLogic ...
type RPCLogic struct {
}

// Register ...
func (rpc *RPCLogic) Register(ctx context.Context, req *proto.RegisterRequest, reply *proto.RegisterReply) (err error) {
	reply.Code = config.FailReplyCode
	u := new(dao.User)
	uData := u.CheckHaveUserName(req.Name)
	if uData.ID > 0 {
		return errors.New("this user name already have , please login")
	}
	u.UserName = req.Name
	u.Password = req.Password
	userID, err := u.Add()
	if err != nil {
		logrus.Infof("register err:%s", err.Error())
		return err
	}
	if userID == 0 {
		return errors.New("register userID empty")
	}

	// set token
	randToken := tools.GetRandomToken(32)
	sessionID := tools.CreateSessionID(randToken)
	userData := make(map[string]interface{})
	userData["userID"] = userID
	userData["userName"] = req.Name
	RedisSessClient.Do("MULTI")
	RedisSessClient.HMSet(sessionID, userData)
	RedisSessClient.Expire(sessionID, 86400*time.Second)
	err = RedisSessClient.Do("EXEC").Err()
	if err != nil {
		logrus.Infof("register set redis token fail!")
		return err
	}
	reply.Code = config.SuccessReplyCode
	reply.AuthToken = randToken
	return
}

// Login ...
func (rpc *RPCLogic) Login(ctx context.Context, args *proto.RegisterRequest, reply *proto.RegisterReply) (err error) {
	reply.Code = config.FailReplyCode
	u := new(dao.User)
	userName := args.Name
	passWord := args.Password
	data := u.CheckHaveUserName(userName)
	if (data.ID == 0) || (passWord != data.Password) {
		return errors.New("no this user or password error")
	}
	loginSessionID := tools.GetSessionIDByUserID(data.ID)
	randToken := tools.GetRandomToken(32)
	sessionID := tools.CreateSessionID(randToken)
	userData := make(map[string]interface{})
	userData["userID"] = data.ID
	userData["userName"] = data.UserName
	// check is login
	token, _ := RedisSessClient.Get(loginSessionID).Result()
	if token != "" {
		// logout already login user session
		oldSession := tools.CreateSessionID(token)
		err := RedisSessClient.Del(oldSession).Err()
		if err != nil {
			return errors.New("logout user fail!token is:" + token)
		}
	}
	RedisSessClient.Do("MULTI")
	RedisSessClient.HMSet(sessionID, userData)
	RedisSessClient.Expire(sessionID, 86400*time.Second)
	RedisSessClient.Set(loginSessionID, randToken, 86400*time.Second)
	err = RedisSessClient.Do("EXEC").Err()
	if err != nil {
		logrus.Infof("register set redis token fail!")
		return err
	}
	reply.Code = config.SuccessReplyCode
	reply.AuthToken = randToken
	return
}

// GetUserInfoByUserID ...
func (rpc *RPCLogic) GetUserInfoByUserID(ctx context.Context, args *proto.GetUserInfoRequest, reply *proto.GetUserInfoResponse) (err error) {
	reply.Code = config.FailReplyCode
	userID := args.UserID
	u := new(dao.User)
	userName := u.GetUserNameByUserID(userID)
	reply.UserID = userID
	reply.UserName = userName
	reply.Code = config.SuccessReplyCode
	return
}

// CheckAuth ...
func (rpc *RPCLogic) CheckAuth(ctx context.Context, args *proto.CheckAuthRequest, reply *proto.CheckAuthResponse) (err error) {
	reply.Code = config.FailReplyCode
	authToken := args.AuthToken
	sessionName := tools.GetSessionName(authToken)
	var userDataMap = map[string]string{}
	userDataMap, err = RedisSessClient.HGetAll(sessionName).Result()
	if err != nil {
		logrus.Infof("check auth fail!,authToken is:%s", authToken)
		return err
	}
	if len(userDataMap) == 0 {
		logrus.Infof("no this user session,authToken is:%s", authToken)
		return
	}
	intUserID, _ := strconv.Atoi(userDataMap["userID"])
	reply.UserID = intUserID
	userName, _ := userDataMap["userName"]
	reply.Code = config.SuccessReplyCode
	reply.UserName = userName
	return
}

// Logout ...
func (rpc *RPCLogic) Logout(ctx context.Context, args *proto.LogoutRequest, reply *proto.LogoutResponse) (err error) {
	reply.Code = config.FailReplyCode
	authToken := args.AuthToken
	sessionName := tools.GetSessionName(authToken)

	var userDataMap = map[string]string{}
	userDataMap, err = RedisSessClient.HGetAll(sessionName).Result()
	if err != nil {
		logrus.Infof("check auth fail!,authToken is:%s", authToken)
		return err
	}
	if len(userDataMap) == 0 {
		logrus.Infof("no this user session,authToken is:%s", authToken)
		return
	}
	intUserID, _ := strconv.Atoi(userDataMap["userID"])
	sessIDMap := tools.GetSessionIDByUserID(intUserID)

	// del sess_map like sess_map_1
	err = RedisSessClient.Del(sessIDMap).Err()
	if err != nil {
		logrus.Infof("logout del sess map error:%s", err.Error())
		return err
	}

	// del serverID
	logic := new(Logic)
	serverIDKey := logic.getUserKey(fmt.Sprintf("%d", intUserID))
	err = RedisSessClient.Del(serverIDKey).Err()
	if err != nil {
		logrus.Infof("logout del server id error:%s", err.Error())
		return err
	}

	err = RedisSessClient.Del(sessionName).Err()
	if err != nil {
		logrus.Infof("logout error:%s", err.Error())
		return err
	}
	reply.Code = config.SuccessReplyCode
	return
}

// Push single send msg
func (rpc *RPCLogic) Push(ctx context.Context, args *proto.Send, reply *proto.SuccessReply) (err error) {
	reply.Code = config.FailReplyCode
	sendData := args
	var bodyBytes []byte
	bodyBytes, err = json.Marshal(sendData)
	if err != nil {
		logrus.Errorf("logic,push msg fail,err:%s", err.Error())
		return
	}

	logic := new(Logic)
	userSidKey := logic.getUserKey(fmt.Sprintf("%d", sendData.ToUserID))
	fmt.Println("userSidkey:", userSidKey)
	serverID := RedisSessClient.Get(userSidKey).Val()
	fmt.Println("serverID:", serverID)
	var serverIDInt int
	serverIDInt, err = strconv.Atoi(serverID)
	if err != nil {
		logrus.Errorf("logic,push parse int fail:%s", err.Error())
		return
	}

	err = logic.RedisPublishChannel(serverIDInt, sendData.ToUserID, bodyBytes)
	if err != nil {
		logrus.Errorf("logic,redis publish err: %s", err.Error())
		return
	}
	reply.Code = config.SuccessReplyCode
	return
}
