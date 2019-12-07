package task

import (
	"encoding/json"
	"gochat/config"
	"gochat/proto"
	"math/rand"

	"github.com/sirupsen/logrus"
)

// PushParams ...
type PushParams struct {
	ServerID int
	UserID   int
	Msg      []byte
	RoomID   int
}

var pushChannel []chan *PushParams

func init() {
	pushChannel = make([]chan *PushParams, config.Conf.Task.TaskBase.PushChan)
}

// GoPush ...
func (task *Task) GoPush() {
	for i := 0; i < len(pushChannel); i++ {
		pushChannel[i] = make(chan *PushParams, config.Conf.Task.TaskBase.PushChanSize)
		go task.processSinglePush(pushChannel[i])
	}
}

func (task *Task) processSinglePush(ch chan *PushParams) {
	var arg *PushParams
	for {
		arg = <-ch
		task.pushSingleToConnect(arg.ServerID, arg.UserID, arg.Msg)
	}
}

// Push ...
func (task *Task) Push(msg string) {
	m := &proto.RedisMsg{}
	if err := json.Unmarshal([]byte(msg), m); err != nil {
		logrus.Infof(" json.Unmarshal err:%v ", err)
	}
	logrus.Infof("push msg info %+v", m)
	switch m.Op {
	case config.OpSingleSend:
		pushChannel[rand.Int()%config.Conf.Task.TaskBase.PushChan] <- &PushParams{
			ServerID: m.ServerID,
			UserID:   m.UserID,
			Msg:      m.Msg,
		}
	case config.OpRoomSend:
		task.broadcastRoomToConnect(m.RoomID, m.Msg)
	case config.OpRoomCountSend:
		task.broadcastRoomCountToConnect(m.RoomID, m.Count)
	case config.OpRoomInfoSend:
		task.broadcastRoomInfoToConnect(m.RoomID, m.RoomUserInfo)
	}
}
