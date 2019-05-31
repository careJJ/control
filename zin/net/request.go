package net

import "zin/zinface"

type Request struct {
	conn zinface.IConnection
	msg zinface.IMessage
}

func NewRequest(conn zinface.IConnection,msg zinface.IMessage) zinface.IRequest {
	req := &Request{
		conn: conn,
		msg:msg,
	}
	return req
}

func (r *Request) GetConnection() zinface.IConnection {
	return r.conn
}
func (r *Request) GetMsg() zinface.IMessage {
	return r.msg
}


