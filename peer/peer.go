package peer

type Client interface {
	GetClientID() uint64
	PostMessage([]byte)
	Close()
	SetUserData(key string, value interface{})
	GetUserData(key string) (interface{}, bool)
	DeleteUserData(key string)
	DeleteAllUserData()
	GetMsgChan() chan []byte
}

type ClientManager interface {
	GetClient(clientID uint64) Client
	GetClientByUserData(key string, v interface{}) Client
	Broadcast(data []byte, condition func(Client) bool)
	CloseAllClient()
}
