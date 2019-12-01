package api

import (
	"context"
	"fmt"
	"gochat/api/router"
	"gochat/api/rpc"
	"gochat/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Chat ...
type Chat struct {
}

// New ...
func New() *Chat {
	return &Chat{}
}

// Run api server,Also, you can use gin,echo ... framework wrap
func (c *Chat) Run() {
	// init logic rpc client
	rpc.InitLogicRPCClient()

	r := router.Register()
	runMode := config.GetGinRunMode()
	logrus.Info("server start , now run mode is ", runMode)
	gin.SetMode(runMode)

	host := config.Conf.API.APIBase.ListenHost
	port := config.Conf.API.APIBase.ListenPort

	// flag.Parse()
	srv := &http.Server{
		Addr:    host + fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("start listen : %s\n", err)
		}
	}()

	// if have two quit signal , this signal will priority capture ,also can graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	logrus.Infof("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Error("Server Shutdown:", err)
	}
	logrus.Infof("Server exiting")
	os.Exit(0)
}
