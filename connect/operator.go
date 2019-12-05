package connect

import "gochat/proto"

// Operator ...
type Operator interface {
	Connect(conn *proto.ConnectRequest) (int, error)
	DisConnect(disConn *proto.DisConnectRequest) (err error)
}

// DefaultOperator ...
type DefaultOperator struct {
}

// Connect rpc call logic layer
func (o *DefaultOperator) Connect(conn *proto.ConnectRequest) (uid int, err error) {
	rpcConnect := new(RPCConnect)
	uid, err = rpcConnect.Connect(conn)
	return
}

// DisConnect rpc call logic layer
func (o *DefaultOperator) DisConnect(disConn *proto.DisConnectRequest) (err error) {
	rpcConnect := new(RPCConnect)
	err = rpcConnect.DisConnect(disConn)
	return
}
