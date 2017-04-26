package dxg

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/zhouyujt/dxg/config"
	"github.com/zhouyujt/dxg/controllers"
	"github.com/zhouyujt/dxg/parsers"
	"github.com/zhouyujt/dxg/peer"
	"github.com/zhouyujt/dxg/routers"
)

var (
	router *routers.Router
	Config *config.Config
)

func init() {
	router = routers.NewRouter()
	Config = config.NewConfig("config.ini")
}

func Run(port int, path string) {
	log.Println("dxg is running...")

	go router.Run(port, path, Config)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	log.Println("dxg stop:", s)
}

func SetRouterPaser(parser parsers.Parser) {
	router.SetParser(parser)
}

func SetRouterRuleFunc(rule func(client peer.Client, routineMgr *routers.LogicRoutineManager) *routers.LogicRoutine) {
	router.SetRuleFunc(rule)
}

func SetRouterRule(rule routers.Rule) {
	router.SetRule(rule)
}

func GetClientManager() peer.ClientManager {
	return router.ClientMgr
}

func RouterHandle(msgID int, c controllers.Controller) {
	router.Handle(msgID, c)
}

func AddRoutine(tickDuration time.Duration, cache ...int) uint64 {
	if len(cache) != 0 {
		return router.RoutineMgr.AddRoutine(tickDuration, cache[0])
	}

	return router.RoutineMgr.AddRoutine(tickDuration)
}

func SetRoutineTicker(routineID uint64, ticker routers.Ticker) {
	r := router.RoutineMgr.GetRoutine(routineID)
	if r != nil {
		r.SetTicker(ticker)
	}
}

func SetRoutineUserData(routineID uint64, key string, v interface{}) {
	r := router.RoutineMgr.GetRoutine(routineID)
	if r != nil {
		r.SetUserData(key, v)
	}
}
