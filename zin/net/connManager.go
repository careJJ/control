package net

import (
	"zin/zinface"
	"sync"
	"fmt"
	"errors"
)

type ConnManager struct {
	connections map[uint32]zinface.IConnection //管理全部的链接
	connLock    sync.RWMutex
}

func NewConnManager() zinface.IConnManager {
	return &ConnManager{
		connections:make(map[uint32] zinface.IConnection),
	}
}
//添加链接
func (connMgr *ConnManager)Add(conn zinface.IConnection){
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	connMgr.connections[conn.GetConnID()]=conn
	fmt.Println("Add connid = ",conn.GetConnID(),"to manager succ!!")
}
//删除链接
func (connMgr *ConnManager) Remove(connID uint32) {
	//加锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connections,  connID)
	fmt.Println("Remove connid = ", connID, " from manager succ!!")

}
//根据链接ID得到链接
func (connMgr *ConnManager) Get(connID uint32) (zinface.IConnection, error) {
	//加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		//找到了
		return conn, nil
	} else {
		//没找到
		return nil, errors.New("connection not FOUND!")
	}
}
//得到目前服务器的链接总个数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}
//清空全部链接的方法
func (connMgr *ConnManager)  ClearConn() {
	//加锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//遍历删除
	for connID, conn := range connMgr.connections {
		//将全部的conn 关闭
		conn.Stop()

		//删除链接
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All Conections succ! conn num = ", connMgr.Len())
}