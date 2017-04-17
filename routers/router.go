package routers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/zhouyujt/dxg/config"
	"github.com/zhouyujt/dxg/controllers"
	"github.com/zhouyujt/dxg/parsers"
	"github.com/zhouyujt/dxg/peer"
	"golang.org/x/net/websocket"
)

const (
	ConnectMsgID    = -1
	DisconnectMsgID = -2
)

type Router struct {
	ClientMgr  *clientManagerImpl
	RoutineMgr *LogicRoutineManager
	webManager *WebManager
	handlers   map[int]controllers.Controller
	parser     parsers.Parser
	ruleFunc   func(client peer.Client, routineMgr *LogicRoutineManager) *LogicRoutine
	rule       Rule
}

type Rule interface {
	Dispathch(client peer.Client, routineMgr *LogicRoutineManager) *LogicRoutine
}

type logicData struct {
	msgID int
	data  []byte
}

func NewRouter() *Router {
	r := new(Router)
	r.ClientMgr = newClientManager()
	r.handlers = make(map[int]controllers.Controller)
	r.RoutineMgr = NewLogicRoutineManager(&r.handlers, r.ClientMgr)
	r.parser = parsers.DefaultParser{}

	return r
}

func (router *Router) SetParser(parser parsers.Parser) {
	router.parser = parser
}

func (router *Router) SetRuleFunc(rule func(client peer.Client, routineMgr *LogicRoutineManager) *LogicRoutine) {
	router.ruleFunc = rule
	router.rule = nil
}

func (router *Router) SetRule(rule Rule) {
	router.rule = rule
	router.ruleFunc = nil
}

// Run accept client request
func (router *Router) Run(port int, path string, cfg *config.Config) {
	router.webManager = NewWebManager(cfg)

	// start websocket server
	go func() {
		if cfg.WebManager.Enable {
			http.Handle("/static/css/", http.FileServer(http.Dir("webmanager")))
			http.Handle("/static/fonts/", http.FileServer(http.Dir("webmanager")))
			http.Handle("/static/img/", http.FileServer(http.Dir("webmanager")))
			http.Handle("/static/js/", http.FileServer(http.Dir("webmanager")))
			http.Handle(cfg.WebManager.Path, router.webManager)
		}
		//http.Handle(path, router)
		http.Handle(path, websocket.Handler(router.OnWebSocket))
		log.Fatal(http.ListenAndServe(`0.0.0.0:`+strconv.Itoa(port), nil))
	}()
}

func (router *Router) OnWebSocket(conn *websocket.Conn) {
	client := newClient(conn, router.parser)
	router.ClientMgr.addClient(client)
	defer router.ClientMgr.delClient(client.GetClientID())
	defer close(client.sendMsgChan)
	defer close(client.closeChan)

	router.dispathMessage(client, ConnectMsgID, make([]byte, 0))

	// message loop
	msgLoopChan := make(chan bool)
	go func() {
		for {
			msgID, msg, disconnect := client.readMessage(0)
			if disconnect {
				router.dispathMessage(client, DisconnectMsgID, make([]byte, 0))
				client.Close()
				break
			}

			router.dispathMessage(client, msgID, msg)
		}
		msgLoopChan <- true
	}()

	// response and push proc
PROC:
	for {
		select {
		case sendMsg, ok := <-client.sendMsgChan:
			if ok {
				client.writeMessage(sendMsg)
			}
		case <-client.closeChan:
			time.Sleep(time.Second) // wait a moment for client.writeMessage
			conn.Close()
			<-msgLoopChan // wait for message loop routine exit
			break PROC
		}
	}
}

func (router *Router) Handle(msgID int, c controllers.Controller) {
	router.handlers[msgID] = c
}

func (router *Router) dispathMessage(client peer.Client, msgID int, data []byte) {
	var routine *LogicRoutine
	if router.ruleFunc != nil {
		routine = router.ruleFunc(client, router.RoutineMgr)

	} else if router.rule != nil {
		routine = router.rule.Dispathch(client, router.RoutineMgr)
	}

	if routine != nil {
		routine.PushQuest(client, msgID, data)
	} else {
		c, ok := router.handlers[msgID]
		if ok {
			c.Proc(client, router.ClientMgr, msgID, data)
		} else {
			j2 := make(map[string]interface{})
			json.Unmarshal(data, &j2)
			log.Println("dispathMessage has no handler!", j2)
		}
	}
}
