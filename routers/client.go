package routers

import (
	"log"
	"sync"

	"github.com/zhouyujt/dxg/parsers"
	"golang.org/x/net/websocket"
)

var nextClientID = uint64(0)
var nextClientIDLocker = sync.Mutex{}

type clientImpl struct {
	id             uint64
	conn           *websocket.Conn
	sendMsgChan    chan []byte
	closeChan      chan bool
	msgParser      parsers.Parser
	userData       map[string]interface{}
	userDataLocker sync.RWMutex
	isClose        bool
	closeLocker    sync.Mutex
	writeLocker    sync.Mutex
}

func newClient(conn *websocket.Conn, msgParser parsers.Parser) *clientImpl {
	c := new(clientImpl)
	c.conn = conn
	c.msgParser = msgParser
	c.id = makeClientID()
	c.sendMsgChan = make(chan []byte)
	c.closeChan = make(chan bool)
	c.userData = make(map[string]interface{})

	return c
}

func (c *clientImpl) readMessage(maxPackageLen uint32) (msgID int, msg []byte, disconnect bool) {
	disconnect = false

	//_, buff, err := c.conn.ReadMessage()
	buff := make([]byte, 0)
	err := websocket.Message.Receive(c.conn, &buff)
	if err != nil {
		log.Println("read message error:", err)
		disconnect = true
		return
	}

	id, contents, err := c.msgParser.Unmarshal(buff, 1)
	if err != nil {
		log.Println(err)
	} else {
		msgID = id
		msg = contents
	}

	return
}

func (c *clientImpl) writeMessage(msg []byte) {
	newMsg := c.msgParser.Marshal(msg, 2)

	c.writeLocker.Lock()
	defer c.writeLocker.Unlock()

	websocket.Message.Send(c.conn, newMsg)
}

func makeClientID() uint64 {
	nextClientIDLocker.Lock()
	defer nextClientIDLocker.Unlock()
	defer func() {
		nextClientID++
	}()

	return nextClientID
}

func (c *clientImpl) GetClientID() uint64 {
	return c.id
}

func (c *clientImpl) PostMessage(msg []byte) {
	c.writeMessage(msg)
}

func (c *clientImpl) Close() {
	c.closeLocker.Lock()
	defer c.closeLocker.Unlock()

	if c.isClose == false {
		c.closeChan <- true
		c.isClose = true
	}
}

func (c *clientImpl) SetUserData(key string, v interface{}) {
	c.userDataLocker.Lock()
	defer c.userDataLocker.Unlock()

	c.userData[key] = v
}

func (c *clientImpl) GetUserData(key string) (v interface{}, ok bool) {
	c.userDataLocker.RLock()
	defer c.userDataLocker.RUnlock()

	v, ok = c.userData[key]

	return
}

func (c *clientImpl) DeleteUserData(key string) {
	c.userDataLocker.Lock()
	defer c.userDataLocker.Unlock()

	delete(c.userData, key)
}

func (c *clientImpl) DeleteAllUserData() {
	c.userDataLocker.Lock()
	defer c.userDataLocker.Unlock()

	c.userData = make(map[string]interface{})
}
