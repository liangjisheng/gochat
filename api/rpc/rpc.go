package rpc

import (
	"context"
	"gochat/config"
	"gochat/proto"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
)

var (
	// LogicRPCClient ...
	LogicRPCClient client.XClient
	once           sync.Once
)

// Logic ...
type Logic struct {
}

// RPCLogicObj ...
var RPCLogicObj *Logic

// InitLogicRPCClient ...
func InitLogicRPCClient() {
	once.Do(func() {
		serviceDiscovery := client.NewEtcdV3Discovery(
			config.Conf.Common.CommonEtcd.BasePath,
			config.Conf.Common.CommonEtcd.ServerPathLogic,
			[]string{config.Conf.Common.CommonEtcd.Host},
			nil,
		)
		LogicRPCClient = client.NewXClient(
			config.Conf.Common.CommonEtcd.ServerPathLogic,
			client.Failtry,
			client.RandomSelect,
			serviceDiscovery,
			client.DefaultOption,
		)

		RPCLogicObj = new(Logic)
	})

	if LogicRPCClient == nil {
		logrus.Errorf("get logic rpc client nil")
	}
}

// Register ...
func (rpc *Logic) Register(req *proto.RegisterRequest) (code int, authToken string, msg string) {
	reply := &proto.RegisterReply{}
	err := LogicRPCClient.Call(context.Background(), "Register", req, reply)
	if err != nil {
		msg = err.Error()
	}
	code = reply.Code
	authToken = reply.AuthToken
	return
}

// Login ...
func (rpc *Logic) Login(req *proto.LoginRequest) (code int, authToken string, msg string) {
	reply := &proto.LoginResponse{}
	logrus.Info("log request", *req)
	err := LogicRPCClient.Call(context.Background(), "Login", req, reply)
	if err != nil {
		msg = err.Error()
	}
	code = reply.Code
	authToken = reply.AuthToken
	return
}

// GetUserNameByUserID ...
func (rpc *Logic) GetUserNameByUserID(req *proto.GetUserInfoRequest) (code int, userName string) {
	reply := &proto.GetUserInfoResponse{}
	LogicRPCClient.Call(context.Background(), "GetUserInfoByUserID", req, reply)
	code = reply.Code
	userName = reply.UserName
	return
}

// CheckAuth ...
func (rpc *Logic) CheckAuth(req *proto.CheckAuthRequest) (code int, userID int, userName string) {
	reply := &proto.CheckAuthResponse{}
	LogicRPCClient.Call(context.Background(), "CheckAuth", req, reply)
	code = reply.Code
	userID = reply.UserID
	userName = reply.UserName
	return
}

// Logout ...
func (rpc *Logic) Logout(req *proto.LogoutRequest) (code int) {
	reply := &proto.LogoutResponse{}
	LogicRPCClient.Call(context.Background(), "Logout", req, reply)
	code = reply.Code
	return
}
