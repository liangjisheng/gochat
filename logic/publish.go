package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gochat/config"
	"gochat/proto"
	"gochat/tools"

	"github.com/go-redis/redis"
	"github.com/rcrowley/go-metrics"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
)

var (
	// RedisClient ...
	RedisClient *redis.Client
	// RedisSessClient ...
	RedisSessClient *redis.Client
)

// InitPublishRedisClient ...
func (logic *Logic) InitPublishRedisClient() (err error) {
	redisOpt := tools.RedisOption{
		Address:  config.Conf.Common.CommonRedis.RedisAddress,
		Password: config.Conf.Common.CommonRedis.RedisPassword,
		Db:       config.Conf.Common.CommonRedis.Db,
	}
	RedisClient = tools.GetRedisInstance(redisOpt)
	if pong, err := RedisClient.Ping().Result(); err != nil {
		logrus.Infof("RedisCli Ping Result pong: %s,  err: %s", pong, err)
	}
	// this can change use another redis save session data
	RedisSessClient = RedisClient
	return err
}

// InitRPCServer ...
func (logic *Logic) InitRPCServer() (err error) {
	var network, addr string
	// a host multi port case
	rpcAddressList := strings.Split(config.Conf.Logic.LogicBase.RPCAddress, ",")
	for _, bind := range rpcAddressList {
		if network, addr, err = tools.ParseNetwork(bind); err != nil {
			logrus.Panicf("InitLogicRpc ParseNetwork error : %s", err.Error())
		}
		logrus.Infof("logic start run at-->%s:%s", network, addr)
		go logic.createRPCServer(network, addr)
	}
	return
}

func (logic *Logic) createRPCServer(network, addr string) {
	s := server.NewServer()
	logic.addRegisterPlugin(s, network, addr)
	// serverID must be unique
	serverPath := config.Conf.Common.CommonEtcd.ServerPathLogic
	serverID := config.Conf.Common.CommonEtcd.ServerID
	err := s.RegisterName(serverPath, new(RPCLogic), fmt.Sprintf("%d", serverID))
	if err != nil {
		logrus.Errorf("register error:%s", err.Error())
	}
	s.RegisterOnShutdown(func(s *server.Server) {
		s.UnregisterAll()
	})
	s.Serve(network, addr)
}

func (logic *Logic) addRegisterPlugin(s *server.Server, network, addr string) {
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

// RedisPublishChannel ...
func (logic *Logic) RedisPublishChannel(serverID int, toUserID int, msg []byte) (err error) {
	redisMsg := proto.RedisMsg{
		Op:       config.OpSingleSend,
		ServerID: serverID,
		UserID:   toUserID,
		Msg:      msg,
	}
	redisMsgStr, err := json.Marshal(redisMsg)
	if err != nil {
		logrus.Errorf("logic,RedisPublishChannel Marshal err:%s", err.Error())
		return err
	}
	redisChannel := config.QueueName
	if err := RedisClient.Publish(redisChannel, redisMsgStr).Err(); err != nil {
		logrus.Errorf("logic,RedisPublishChannel err:%s", err.Error())
		return err
	}
	return
}

func (logic *Logic) getRoomUserKey(authKey string) string {
	var returnKey bytes.Buffer
	returnKey.WriteString(config.RedisRoomPrefix)
	returnKey.WriteString(authKey)
	return returnKey.String()
}

func (logic *Logic) getRoomOnlineCountKey(authKey string) string {
	var returnKey bytes.Buffer
	returnKey.WriteString(config.RedisRoomOnlinePrefix)
	returnKey.WriteString(authKey)
	return returnKey.String()
}

func (logic *Logic) getUserKey(authKey string) string {
	var returnKey bytes.Buffer
	returnKey.WriteString(config.RedisPrefix)
	returnKey.WriteString(authKey)
	return returnKey.String()
}
