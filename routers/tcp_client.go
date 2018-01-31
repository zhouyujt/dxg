package routers

import (
	"log"
	"net"
	"sync"

	"github.com/zhouyujt/dxg/parsers"
)

type tcpClientImpl struct {
	baseClientImpl
	conn       net.Conn
	readLocker sync.Mutex
}

func newTcpClient(conn net.Conn, msgParser parsers.Parser) *tcpClientImpl {
	c := new(tcpClientImpl)
	c.conn = conn
	c.msgParser = msgParser
	c.id = makeClientID()
	c.sendMsgChan = make(chan []byte)
	c.closeChan = make(chan bool)
	c.userData = make(map[string]interface{})

	return c
}

func (c *tcpClientImpl) readMessage(maxPackageLen uint32) (msgID int, msg []byte, disconnect bool) {
	c.readLocker.Lock()
	defer c.readLocker.Unlock()
	disconnect = false

	buff := make([]byte, 0)
	_, err := c.conn.Read(buff)
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

func (c *tcpClientImpl) writeMessage(msg []byte) {
	newMsg := c.msgParser.Marshal(msg, 2)

	c.writeLocker.Lock()
	defer c.writeLocker.Unlock()

	c.conn.Write(newMsg)
}

func (c *tcpClientImpl) PostMessage(msg []byte) {
	c.writeMessage(msg)
}
