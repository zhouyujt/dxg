package routers

import (
	"log"

	"github.com/zhouyujt/dxg/parsers"
	"golang.org/x/net/websocket"
)

type wsClientImpl struct {
	baseClientImpl
	conn *websocket.Conn
}

func newWsClient(conn *websocket.Conn, msgParser parsers.Parser) *wsClientImpl {
	c := new(wsClientImpl)
	c.conn = conn
	c.msgParser = msgParser
	c.id = makeClientID()
	c.sendMsgChan = make(chan []byte)
	c.closeChan = make(chan bool)
	c.userData = make(map[string]interface{})

	return c
}

func (c *wsClientImpl) readMessage(maxPackageLen uint32) (msgID int, msg []byte, disconnect bool) {
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

func (c *wsClientImpl) writeMessage(msg []byte) {
	newMsg := c.msgParser.Marshal(msg, 2)

	c.writeLocker.Lock()
	defer c.writeLocker.Unlock()

	websocket.Message.Send(c.conn, newMsg)
}

func (c *wsClientImpl) PostMessage(msg []byte) {
	c.writeMessage(msg)
}
