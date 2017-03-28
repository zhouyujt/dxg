package controllers

import "github.com/zhouyujt/dxg/peer"

type Controller interface {
	Proc(c peer.Client, cm peer.ClientManager, msgID int, data []byte)
}
