package routers

import (
	"sync"

	"github.com/zhouyujt/dxg/parsers"
)

var nextClientID = uint64(0)
var nextClientIDLocker = sync.Mutex{}

type baseClientImpl struct {
	id             uint64
	sendMsgChan    chan []byte
	closeChan      chan bool
	msgParser      parsers.Parser
	userData       map[string]interface{}
	userDataLocker sync.RWMutex
	isClose        bool
	closeLocker    sync.Mutex
	writeLocker    sync.Mutex
}

func makeClientID() uint64 {
	nextClientIDLocker.Lock()
	defer nextClientIDLocker.Unlock()
	defer func() {
		nextClientID++
	}()

	return nextClientID
}

func (c *baseClientImpl) GetClientID() uint64 {
	return c.id
}

func (c *baseClientImpl) Close() {
	c.closeLocker.Lock()
	defer c.closeLocker.Unlock()

	if c.isClose == false {
		c.closeChan <- true
		c.isClose = true
	}
}

func (c *baseClientImpl) SetUserData(key string, v interface{}) {
	c.userDataLocker.Lock()
	defer c.userDataLocker.Unlock()

	c.userData[key] = v
}

func (c *baseClientImpl) GetUserData(key string) (v interface{}, ok bool) {
	c.userDataLocker.RLock()
	defer c.userDataLocker.RUnlock()

	v, ok = c.userData[key]

	return
}

func (c *baseClientImpl) DeleteUserData(key string) {
	c.userDataLocker.Lock()
	defer c.userDataLocker.Unlock()

	delete(c.userData, key)
}

func (c *baseClientImpl) DeleteAllUserData() {
	c.userDataLocker.Lock()
	defer c.userDataLocker.Unlock()

	c.userData = make(map[string]interface{})
}

func (c *baseClientImpl) GetMsgChan() chan []byte {
	return c.sendMsgChan
}
