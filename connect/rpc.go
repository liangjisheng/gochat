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
