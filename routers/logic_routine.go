package routers

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/zhouyujt/dxg/controllers"
	"github.com/zhouyujt/dxg/peer"
)

var nextRoutineID = uint64(0)
var nextRoutineIDLocker = sync.Mutex{}

type LogicRoutine struct {
	id             uint64
	requestChan    chan requestData
	tickDuration   time.Duration
	userData       map[string]interface{}
	userDataLocker sync.RWMutex
	handlers       *map[int]controllers.Controller
	clientMgr      peer.ClientManager
	ticker         Ticker
	tickerLocker   sync.Mutex
}

type Ticker interface {
	Tick(tickDuration time.Duration)
}

type requestData struct {
	client peer.Client
	msgID  int
	data   []byte
}

func NewLogicRoutine(tickDuration time.Duration, handlers *map[int]controllers.Controller, clientMgr peer.ClientManager, cache int) *LogicRoutine {
	r := new(LogicRoutine)
	r.requestChan = make(chan requestData, cache)
	r.tickDuration = tickDuration
	r.id = makeRoutineID()
	r.userData = make(map[string]interface{})
	r.handlers = handlers
	r.clientMgr = clientMgr

	go r.heartbeat()
	return r
}

func (routine *LogicRoutine) SetTicker(ticker Ticker) {
	routine.tickerLocker.Lock()
	defer routine.tickerLocker.Unlock()

	routine.ticker = ticker
}

func (routine *LogicRoutine) GetID() uint64 {
	return routine.id
}

func (routine *LogicRoutine) PushQuest(client peer.Client, msgID int, data []byte) {
	routine.requestChan <- requestData{client: client, msgID: msgID, data: data}
}

func (routine *LogicRoutine) SetUserData(key string, v interface{}) {
	routine.userDataLocker.Lock()
	defer routine.userDataLocker.Unlock()

	routine.userData[key] = v
}

func (routine *LogicRoutine) GetUserData(key string) (v interface{}, ok bool) {
	routine.userDataLocker.RLock()
	routine.userDataLocker.RUnlock()

	v, ok = routine.userData[key]

	return
}

func makeRoutineID() uint64 {
	nextRoutineIDLocker.Lock()
	defer nextRoutineIDLocker.Unlock()
	defer func() {
		nextRoutineID++
	}()

	return nextRoutineID
}

func (routine *LogicRoutine) heartbeat() {
	t := time.NewTicker(routine.tickDuration)
	begin := time.Now()
	for {
		routine.tickerLocker.Lock()
		if routine.ticker != nil {
			routine.ticker.Tick(routine.tickDuration)
		}
		routine.tickerLocker.Unlock()
	BREAK:
		for {
			select {
			case data, ok := <-routine.requestChan:
				if ok {
					c, ok := (*routine.handlers)[data.msgID]
					if ok {
						c.Proc(data.client, routine.clientMgr, data.msgID, data.data)
					} else {
						tmp := make(map[string]interface{})
						json.Unmarshal(data.data, &tmp)
						log.Println("dxg:dispathMessage has no handler!", tmp)
					}
				}
			case end := <-t.C:
				tt := end.Sub(begin)
				if tt-time.Millisecond*100 > routine.tickDuration {
					log.Println("dxg:logic routine proc is too busy!!!:", routine.GetID(), tt)
				}
				begin = end
				break BREAK
			}
		}
	}
}
