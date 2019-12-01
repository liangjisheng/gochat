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
