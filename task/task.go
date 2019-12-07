package task

import (
	"gochat/config"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Task ...
type Task struct {
}

// New ...
func New() *Task {
	return new(Task)
}

// Run ...
func (task *Task) Run() {
	//read config
	taskConfig := config.Conf.Task
	runtime.GOMAXPROCS(taskConfig.TaskBase.CPUNum)
	// read from redis queue
	if err := task.InitSubscribeRedisClient(); err != nil {
		logrus.Panicf("task init publishRedisClient fail,err:%s", err.Error())
	}
	// rpc call connect layer send msg
	if err := task.InitConnectRPCClient(); err != nil {
		logrus.Panicf("task init InitConnectRpcClient fail,err:%s", err.Error())
	}
	// GoPush
	task.GoPush()
}
