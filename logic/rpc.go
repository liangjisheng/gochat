package logic

import (
	"context"
	"errors"
	"gochat/tools"
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
