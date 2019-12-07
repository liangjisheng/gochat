package task

import (
	"context"
	"encoding/json"
	"gochat/config"
	"gochat/proto"
	"gochat/tools"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
)

// RPCConnectClientList ...
var RPCConnectClientList map[int]client.XClient

// InitConnectRPCClient ...
func (task *Task) InitConnectRPCClient() (err error) {
	etcdConfig := config.Conf.Common.CommonEtcd
	d := client.NewEtcdV3Discovery(etcdConfig.BasePath,
		etcdConfig.ServerPathConnect,
		[]string{etcdConfig.Host}, nil)
	if len(d.GetServices()) <= 0 {
		logrus.Panicf("no etcd server find!")
	}

	RPCConnectClientList = make(map[int]client.XClient, len(d.GetServices()))
	for _, connectConf := range d.GetServices() {
		connectConf.Value = strings.Replace(connectConf.Value, "=&tps=0", "", 1)
		serverID, err := strconv.ParseInt(connectConf.Value, 10, 8)
		if err != nil {
			logrus.Panicf("InitComets err，Can't find serverId. error: %s", err)
		}
		d := client.NewPeer2PeerDiscovery(connectConf.Key, "")
		RPCConnectClientList[int(serverID)] = client.NewXClient(
			etcdConfig.ServerPathConnect, client.Failtry,
			client.RandomSelect, d, client.DefaultOption)
		logrus.Infof("InitConnectRpcClient addr %s, v %+v", connectConf.Key, RPCConnectClientList[int(serverID)])
	}
	return
}

func (task *Task) pushSingleToConnect(serverID int, userID int, msg []byte) {
	logrus.Infof("pushSingleToConnect Body %s", string(msg))
	pushMsgReq := &proto.PushMsgRequest{
		UserID: userID,
		Msg: proto.Msg{
			Ver:       config.MsgVersion,
			Operation: config.OpSingleSend,
			SeqID:     tools.GetSnowflakeID(),
			Body:      msg,
		},
	}
	reply := &proto.SuccessReply{}
	// todo lock
	err := RPCConnectClientList[serverID].Call(context.Background(), "PushSingleMsg", pushMsgReq, reply)
	if err != nil {
		logrus.Infof(" pushSingleToConnect Call err %v", err)
	}
	logrus.Infof("reply %s", reply.Msg)
}

func (task *Task) broadcastRoomToConnect(roomID int, msg []byte) {
	pushRoomMsgReq := &proto.PushRoomMsgRequest{
		RoomID: roomID,
		Msg: proto.Msg{
			Ver:       config.MsgVersion,
			Operation: config.OpRoomSend,
			SeqID:     tools.GetSnowflakeID(),
			Body:      msg,
		},
	}
	reply := &proto.SuccessReply{}
	for _, rpc := range RPCConnectClientList {
		logrus.Infof("broadcastRoomToConnect rpc  %v", rpc)
		rpc.Call(context.Background(), "PushRoomMsg", pushRoomMsgReq, reply)
		logrus.Infof("reply %s", reply.Msg)
	}
}

func (task *Task) broadcastRoomCountToConnect(roomID, count int) {
	msg := &proto.RedisRoomCountMsg{
		Count: count,
		Op:    config.OpRoomCountSend,
	}
	var body []byte
	var err error
	if body, err = json.Marshal(msg); err != nil {
		logrus.Warnf("broadcastRoomCountToConnect  json.Marshal err :%s", err.Error())
		return
	}
	pushRoomMsgReq := &proto.PushRoomMsgRequest{
		RoomID: roomID,
		Msg: proto.Msg{
			Ver:       config.MsgVersion,
			Operation: config.OpRoomCountSend,
			SeqID:     tools.GetSnowflakeID(),
			Body:      body,
		},
	}
	reply := &proto.SuccessReply{}
	for _, rpc := range RPCConnectClientList {
		logrus.Infof("broadcastRoomCountToConnect rpc  %v", rpc)
		rpc.Call(context.Background(), "PushRoomCount", pushRoomMsgReq, reply)
		logrus.Infof("reply %s", reply.Msg)
	}
}

func (task *Task) broadcastRoomInfoToConnect(roomID int, roomUserInfo map[string]string) {
	msg := &proto.RedisRoomInfo{
		Count:        len(roomUserInfo),
		Op:           config.OpRoomInfoSend,
		RoomUserInfo: roomUserInfo,
		RoomID:       roomID,
	}
	var body []byte
	var err error
	if body, err = json.Marshal(msg); err != nil {
		logrus.Warnf("broadcastRoomInfoToConnect  json.Marshal err :%s", err.Error())
		return
	}
	pushRoomMsgReq := &proto.PushRoomMsgRequest{
		RoomID: roomID,
		Msg: proto.Msg{
			Ver:       config.MsgVersion,
			Operation: config.OpRoomInfoSend,
			SeqID:     tools.GetSnowflakeID(),
			Body:      body,
		},
	}
	reply := &proto.SuccessReply{}
	for _, rpc := range RPCConnectClientList {
		logrus.Infof("broadcastRoomInfoToConnect rpc  %v", rpc)
		rpc.Call(context.Background(), "PushRoomInfo", pushRoomMsgReq, reply)
		logrus.Infof("broadcastRoomInfoToConnect rpc  reply %v", reply)
	}
}
