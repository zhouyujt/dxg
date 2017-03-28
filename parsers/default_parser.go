package parsers

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"errors"
	"fmt"
)

const (
	PacketHeadLen uint32 = 2 /*proto flag*/ + 4 /*packet len*/ + 4 /*reserved*/
	MaxPackageLen        = 0
)

type DefaultParser struct {
}

type Packet struct {
	MsgID    int `json:"msg_id"`
	Contents []byte
}

// Unmarshal parser the data.
//              +---------------------------------------------------------------------------------------------------------------------------------------+
// data format: |   0x01 (1 bytes)  |   request(0x01) or respond(0x02) (1 bytes)    |   packet len (4 bytes)    |   reserved (4 bytes)  |   contents    |
//              +---------------------------------------------------------------------------------------------------------------------------------------+
func (parser DefaultParser) Unmarshal(data []byte, magic int) (msgID int, contents []byte, err error) {
	packet := Packet{}
	if 0 == bytes.Compare(data[:2], []byte{0x01, byte(magic)}) {
		var packageLen uint32
		var reserved uint32
		binary.Read(bytes.NewBuffer(data[2:6]), binary.BigEndian, &packageLen)
		binary.Read(bytes.NewBuffer(data[6:10]), binary.BigEndian, &reserved)
		if 0 == MaxPackageLen || packageLen < MaxPackageLen {
			message := data[PacketHeadLen:]
			if int(packageLen) == len(message) {
				err = json.Unmarshal(message, &packet)
				if err != nil {
					err = errors.New("DefaultParsor Unmarshal failed")
				} else {
					contents = message
				}
			}
		} else {
			tmp := fmt.Sprintf("DefaultParsor over MaxPackageLen: %d,drop whole package", packageLen)
			err = errors.New(tmp)
		}
	}

	msgID = packet.MsgID

	return
}

func (parser DefaultParser) Marshal(msg []byte, magic int) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, []byte{0x01, byte(magic)})
	binary.Write(buf, binary.BigEndian, uint32(len(msg)))
	binary.Write(buf, binary.BigEndian, []byte{0, 0, 0, 0})
	binary.Write(buf, binary.BigEndian, msg)

	return buf.Bytes()
}
