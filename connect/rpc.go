package connect

import (
	"context"
	"errors"
	"fmt"
	"gochat/config"
	"gochat/proto"
	"gochat/tools"
	"strings"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
)

var logicRPCClient client.XClient
var once sync.Once

// RPCConnect ...
type RPCConnect struct {
}

// InitLogicRPCClient ...
func (c *Connect) InitLogicRPCClient() (err error) {
	once.Do(func() {
		d := client.NewEtcdV3Discovery(
			config.Conf.Common.CommonEtcd.BasePath,
			config.Conf.Common.CommonEtcd.ServerPathLogic,
			[]string{config.Conf.Common.CommonEtcd.Host},
			nil,
		)
		logicRPCClient = client.NewXClient(
			config.Conf.Common.CommonEtcd.ServerPathLogic,
			client.Failtry,
			client.RandomSelect,
			d,
			client.DefaultOption,
		)
	})
	if logicRPCClient == nil {
		return errors.New("get rpc client nil")
	}
	return
}

// InitConnectRPCServer ...
func (c *Connect) InitConnectRPCServer() (err error) {
	var network, addr string
	connectRPCAddress := strings.Split(config.Conf.Connect.ConnectRPCAddress.Address, ",")
	for _, bind := range connectRPCAddress {
		if network, addr, err = tools.ParseNetwork(bind); err != nil {
			logrus.Panicf("InitConnectRpcServer ParseNetwork error : %s", err)
		}
		logrus.Infof("Connect start run at-->%s:%s", network, addr)
		go c.createConnectRPCServer(network, addr)
	}
	return
}

func (c *Connect) createConnectRPCServer(network string, addr string) {
	s := server.NewServer()
	addRegistryPlugin(s, network, addr)
	s.RegisterName(
		config.Conf.Common.CommonEtcd.ServerPathConnect,
		new(RPCConnectPush),
		fmt.Sprintf("%d", config.Conf.Common.CommonEtcd.ServerID),
	)
	s.Serve(network, addr)
}

func addRegistryPlugin(s *server.Server, network string, addr string) {
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: network + "@" + addr,
		EtcdServers:    []string{config.Conf.Common.CommonEtcd.Host},
		BasePath:       config.Conf.Common.CommonEtcd.BasePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		logrus.Fatal(err)
	}
	s.Plugins.Add(r)
}

// Connect ...
func (rpc *RPCConnect) Connect(connReq *proto.ConnectRequest) (uid int, err error) {
	reply := &proto.ConnectReply{}
	err = logicRPCClient.Call(context.Background(), "Connect", connReq, reply)
	if err != nil {
		logrus.Fatalf("failed to call: %v", err)
	}
	uid = reply.UserID
	logrus.Infof("connect logic userId :%d", reply.UserID)
	return
}

// DisConnect ...
func (rpc *RPCConnect) DisConnect(disConnReq *proto.DisConnectRequest) (err error) {
	reply := &proto.DisConnectReply{}
	if err = logicRPCClient.Call(context.Background(), "DisConnect", disConnReq, reply); err != nil {
		logrus.Fatalf("failed to call: %v", err)
	}
	return
}

// RPCConnectPush ...
type RPCConnectPush struct {
}

// PushSingleMsg ...
func (rpc *RPCConnectPush) PushSingleMsg(ctx context.Context, pushMsgReq *proto.PushMsgRequest, successReply *proto.SuccessReply) (err error) {
	var (
		bucket  *Bucket
		channel *Channel
	)
	logrus.Infof("rpc PushMsg :%v ", pushMsgReq)
	if pushMsgReq == nil {
		logrus.Errorf("rpc PushSingleMsg() args:(%v)", pushMsgReq)
		return
	}
	bucket = DefaultServer.Bucket(pushMsgReq.UserID)
	if channel = bucket.Channel(pushMsgReq.UserID); channel != nil {
		err = channel.Push(&pushMsgReq.Msg)
		logrus.Infof("DefaultServer Channel err nil ,args: %v", pushMsgReq)
		return
	}
	successReply.Code = config.SuccessReplyCode
	successReply.Msg = config.SuccessReplyMsg
	logrus.Infof("successReply:%v", successReply)
	return
}

// PushRoomMsg ...
func (rpc *RPCConnectPush) PushRoomMsg(ctx context.Context, pushRoomMsgReq *proto.PushRoomMsgRequest, successReply *proto.SuccessReply) (err error) {
	successReply.Code = config.SuccessReplyCode
	successReply.Msg = config.SuccessReplyMsg
	logrus.Infof("PushRoomMsg msg %+v", pushRoomMsgReq)
	for _, bucket := range DefaultServer.Buckets {
		bucket.BroadcastRoom(pushRoomMsgReq)
	}
	return
}

// PushRoomCount ...
func (rpc *RPCConnectPush) PushRoomCount(ctx context.Context, pushRoomMsgReq *proto.PushRoomMsgRequest, successReply *proto.SuccessReply) (err error) {
	successReply.Code = config.SuccessReplyCode
	successReply.Msg = config.SuccessReplyMsg
	logrus.Infof("PushRoomCount msg %v", pushRoomMsgReq)
	for _, bucket := range DefaultServer.Buckets {
		bucket.BroadcastRoom(pushRoomMsgReq)
	}
	return
}

// PushRoomInfo ...
func (rpc *RPCConnectPush) PushRoomInfo(ctx context.Context, pushRoomMsgReq *proto.PushRoomMsgRequest, successReply *proto.SuccessReply) (err error) {
	successReply.Code = config.SuccessReplyCode
	successReply.Msg = config.SuccessReplyMsg
	logrus.Infof("connect,PushRoomInfo msg %+v", pushRoomMsgReq)
	for _, bucket := range DefaultServer.Buckets {
		bucket.BroadcastRoom(pushRoomMsgReq)
	}
	return
}
