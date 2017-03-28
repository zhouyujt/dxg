package dxg

import (
	"log"
	"os"
	"os/signal"

	"github.com/zhouyujt/dxg/config"
	"github.com/zhouyujt/dxg/routers"
)

var (
	Router *routers.Router
	Config *config.Config
)

func init() {
	Router = routers.NewRouter()
	Config = config.NewConfig("config.ini")
}

func Run(port int, path string) {
	log.Println("dxg is running...")

	go Router.Run(port, path, Config)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	log.Println("dxg stop:", s)
}
