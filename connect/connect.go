package connect

import (
	"gochat/config"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

// DefaultServer ...
var DefaultServer *Server

// Connect ...
type Connect struct {
}

// New ...
func New() *Connect {
	return new(Connect)
}

// Run ...
func (c *Connect) Run() {
	// get Connect layer config
	connectConfig := config.Conf.Connect

	// set the maximum number of CPUs that can be executing
	runtime.GOMAXPROCS(connectConfig.ConnectBucket.CPUNum)

	// init logic layer rpc client, call logic layer rpc server
	if err := c.InitLogicRPCClient(); err != nil {
		logrus.Panicf("InitLogicRpcClient err:%s", err.Error())
	}

	// init Connect layer rpc server, logic client will call this
	Buckets := make([]*Bucket, connectConfig.ConnectBucket.CPUNum)
	for i := 0; i < connectConfig.ConnectBucket.CPUNum; i++ {
		Buckets[i] = NewBucket(BucketOptions{
			ChannelSize:   connectConfig.ConnectBucket.Channel,
			RoomSize:      connectConfig.ConnectBucket.Room,
			RoutineAmount: connectConfig.ConnectBucket.RoutineAmount,
			RoutineSize:   connectConfig.ConnectBucket.RoutineSize,
		})
	}
	operator := new(DefaultOperator)
	DefaultServer = NewServer(Buckets, operator, ServerOptions{
		WriteWait:       10 * time.Second,
		PongWait:        60 * time.Second,
		PingPeriod:      54 * time.Second,
		MaxMessageSize:  512,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		BroadcastSize:   512,
	})

	// init Connect layer rpc server ,task layer will call this
	if err := c.InitConnectRPCServer(); err != nil {
		logrus.Panicf("InitConnectRpcServer Fatal error: %s \n", err)
	}

	// start Connect layer server handler persistent connection
	if err := c.InitWebsocket(); err != nil {
		logrus.Panicf("Connect layer InitWebsocket() error:  %s \n", err.Error())
	}
}
