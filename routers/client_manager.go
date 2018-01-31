package routers

import (
	"log"
	"sync"

	"github.com/zhouyujt/dxg/peer"
)

type clientManagerImpl struct {
	clientLocker sync.RWMutex
	clientMap    map[uint64]peer.Client
}

func newClientManager() *clientManagerImpl {
	m := new(clientManagerImpl)
	m.clientMap = make(map[uint64]peer.Client)

	return m
}

func (m *clientManagerImpl) addClient(c peer.Client) {
	m.clientLocker.Lock()
	defer m.clientLocker.Unlock()

	m.clientMap[c.GetClientID()] = c
}

func (m *clientManagerImpl) delClient(clientID uint64) {
	m.clientLocker.Lock()
	defer m.clientLocker.Unlock()

	delete(m.clientMap, clientID)
}

func (m *clientManagerImpl) GetClient(clientID uint64) peer.Client {
	m.clientLocker.RLock()
	defer m.clientLocker.RUnlock()

	var c peer.Client
	c, _ = m.clientMap[clientID]

	return c
}

func (m *clientManagerImpl) GetClientByUserData(key string, v interface{}) peer.Client {
	m.clientLocker.RLock()
	defer m.clientLocker.RUnlock()

	var c peer.Client
	for _, test := range m.clientMap {
		data, ok := test.GetUserData(key)
		if ok {
			if data == v {
				c = test
				break
			}
		}
	}

	return c
}

func (m *clientManagerImpl) Broadcast(data []byte, condition func(peer.Client) bool) {
	m.clientLocker.RLock()
	defer m.clientLocker.RUnlock()

	for _, c := range m.clientMap {
		if condition(c) {
			go func(pushMsgChan chan []byte) {
				defer func() {
					if err := recover(); err != nil {
						log.Println("broadcastMessage error:", err)
					}
				}()
				pushMsgChan <- data
			}(c.GetMsgChan())
		}
	}
}

func (m *clientManagerImpl) CloseAllClient() {
	m.clientLocker.RLock()
	defer m.clientLocker.RUnlock()

	for _, c := range m.clientMap {
		c.Close()
	}
}
