package logic

import (
	"gochat/config"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Logic ...
type Logic struct {
}

// New ...
func New() *Logic {
	return new(Logic)
}

// Run ...
func (logic *Logic) Run() {
	// read config
	logicConfig := config.Conf.Logic
	runtime.GOMAXPROCS(logicConfig.LogicBase.CPUNum)

	// init publish redis
	if err := logic.InitPublishRedisClient(); err != nil {
		logrus.Panicf("logic init publishRedisClient fail,err:%s", err.Error())
	}

	// init rpc server
	if err := logic.InitRPCServer(); err != nil {
		logrus.Panicf("logic init rpc server fail")
	}
}
