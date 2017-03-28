package routers

import (
	"sync"
	"time"

	"github.com/zhouyujt/dxg/controllers"
	"github.com/zhouyujt/dxg/peer"
)

type LogicRoutineManager struct {
	routines       map[uint64]*LogicRoutine
	routinesLocker sync.RWMutex
	handlers       *map[int]controllers.Controller
	clientMgr      peer.ClientManager
}

func NewLogicRoutineManager(handlers *map[int]controllers.Controller, clientMgr peer.ClientManager) *LogicRoutineManager {
	mgr := new(LogicRoutineManager)
	mgr.routines = make(map[uint64]*LogicRoutine)
	mgr.handlers = handlers
	mgr.clientMgr = clientMgr
	return mgr
}

// AddRoutine returns a new LogicRoutine's id.the new routine will working when router's rule dispached.
func (mgr *LogicRoutineManager) AddRoutine(tickDuration time.Duration) uint64 {
	mgr.routinesLocker.Lock()
	defer mgr.routinesLocker.Unlock()

	r := NewLogicRoutine(tickDuration, mgr.handlers, mgr.clientMgr)
	id := r.GetID()
	mgr.routines[id] = r

	return id
}

// GetRoutine returns a pointer to LogicRoutine by id.
func (mgr *LogicRoutineManager) GetRoutine(id uint64) *LogicRoutine {
	mgr.routinesLocker.RLock()
	defer mgr.routinesLocker.RUnlock()

	r, _ := mgr.routines[id]

	return r
}

func (mgr *LogicRoutineManager) GetRoutineByUserData(key string, v interface{}) *LogicRoutine {
	mgr.routinesLocker.RLock()
	defer mgr.routinesLocker.RUnlock()

	var routine *LogicRoutine
	for _, test := range mgr.routines {
		data, ok := test.GetUserData(key)
		if ok {
			if data == v {
				routine = test
				break
			}
		}
	}

	return routine
}
