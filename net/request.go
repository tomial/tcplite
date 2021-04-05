package net

import "github.com/tomial/tcplite/iface"

type Request struct {
	// 连接对象
	Conn iface.IConnection

	// 数据
	msg iface.IMessage
}

func (r *Request) GetConnection() iface.IConnection {
	return r.Conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
