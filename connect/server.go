package connect

import (
	"encoding/json"
	"fmt"
	"gochat/config"
	"gochat/proto"
	"gochat/tools"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Server ...
type Server struct {
	Buckets   []*Bucket
	Options   ServerOptions
	bucketIDx uint32
	operator  Operator
}

// ServerOptions ...
type ServerOptions struct {
	WriteWait       time.Duration
	PongWait        time.Duration
	PingPeriod      time.Duration
	MaxMessageSize  int64
	ReadBufferSize  int
	WriteBufferSize int
	BroadcastSize   int
}

// NewServer ...
func NewServer(b []*Bucket, o Operator, options ServerOptions) *Server {
	s := new(Server)
	s.Buckets = b
	s.Options = options
	s.bucketIDx = uint32(len(b))
	s.operator = o
	return s
}

// Bucket reduce lock competition, use google city hash insert to different bucket
func (s *Server) Bucket(userID int) *Bucket {
	userIDStr := fmt.Sprintf("%d", userID)
	idx := tools.CityHash32([]byte(userIDStr), uint32(len(userIDStr))) % s.bucketIDx
	return s.Buckets[idx]
}

// send data to websocket conn
func (s *Server) writePump(ch *Channel) {
	// PingPeriod default eq 54s
	ticker := time.NewTicker(s.Options.PingPeriod)
	defer func() {
		ticker.Stop()
		ch.conn.Close()
	}()

	for {
		select {
		case message, ok := <-ch.broadcast:
			// write data dead time , like http timeout , default 10s
			ch.conn.SetWriteDeadline(time.Now().Add(s.Options.WriteWait))
			if !ok {
				logrus.Warn("SetWriteDeadline not ok")
				ch.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := ch.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logrus.Warnf(" ch.conn.NextWriter err :%s  ", err.Error())
				return
			}
			logrus.Infof("message write body:%s", message.Body)
			w.Write(message.Body)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// heartbeat，if ping error will exit and close current websocket conn
			ch.conn.SetWriteDeadline(time.Now().Add(s.Options.WriteWait))
			logrus.Infof("websocket.PingMessage :%v", websocket.PingMessage)
			if err := ch.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// get data from websocket conn
func (s *Server) readPump(ch *Channel) {
	defer func() {
		logrus.Infof("start exec disConnect ...")
		if ch.Room == nil || ch.userID == 0 {
			logrus.Infof("roomId and userID eq 0")
			ch.conn.Close()
			return
		}
		logrus.Infof("exec disConnect ...")
		disConnectRequest := new(proto.DisConnectRequest)
		disConnectRequest.RoomID = ch.Room.ID
		disConnectRequest.UserID = ch.userID
		s.Bucket(ch.userID).DeleteChannel(ch)
		if err := s.operator.DisConnect(disConnectRequest); err != nil {
			logrus.Warnf("DisConnect err :%s", err.Error())
		}
		ch.conn.Close()
	}()

	ch.conn.SetReadLimit(s.Options.MaxMessageSize)
	ch.conn.SetReadDeadline(time.Now().Add(s.Options.PongWait))
	ch.conn.SetPongHandler(func(string) error {
		ch.conn.SetReadDeadline(time.Now().Add(s.Options.PongWait))
		return nil
	})

	for {
		_, message, err := ch.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Errorf("readPump ReadMessage err:%s", err.Error())
				return
			}
		}
		if message == nil {
			return
		}

		var connReq *proto.ConnectRequest
		logrus.Infof("get a message :%s", message)
		if err := json.Unmarshal([]byte(message), &connReq); err != nil {
			logrus.Errorf("message struct %+v", connReq)
		}
		if connReq.AuthToken == "" {
			logrus.Errorf("s.operator.Connect no authToken")
			return
		}

		connReq.ServerID = config.Conf.Connect.ConnectBase.ServerID
		userID, err := s.operator.Connect(connReq)
		if err != nil {
			logrus.Errorf("s.operator.Connect error %s", err.Error())
			return
		}
		if userID == 0 {
			logrus.Error("Invalid AuthToken ,userID empty")
			return
		}

		logrus.Infof("websocket rpc call return userId:%d,RoomId:%d", userID, connReq.RoomID)
		b := s.Bucket(userID)
		// insert into a bucket
		err = b.Put(userID, connReq.RoomID, ch)
		if err != nil {
			logrus.Errorf("conn close err: %s", err.Error())
			ch.conn.Close()
		}
	}
}
