package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gochat/api"
	"gochat/logic"
)

func main() {
	var module string
	flag.StringVar(&module, "module", "", "assign run module")
	flag.Parse()
	fmt.Println(fmt.Sprintf("start run %s module", module))
	switch module {
	case "logic":
		logic.New().Run()
	case "api":
		api.New().Run()
	default:
		fmt.Println("exiting,module param error!")
		return
	}

	fmt.Println(fmt.Sprintf("run %s module done!", module))
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	fmt.Println("Server exiting")
}
