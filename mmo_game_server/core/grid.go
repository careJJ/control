package core

import (
	"sync"
	"fmt"
)

type Grid struct{
	//格子ID
	GID int
	//格子的左边边界的坐标
	MinX int
	//格子的右边边界的坐标
	MaxX int
	//格子的上边边界的坐标
	MinY int
	//格子的下边边界的坐标
	MaxY int
	//当前格子内玩家/物体 成员的ID几个 map[玩家/物体ID]
	playerIDs map[int]interface{}
	//保护当前格子本荣的map的锁
	pIDLock sync.RWMutex

}

//初始化格子的方法
func NewGrid(gID,minX,maxX,minY,maxY int) *Grid{
	return &Grid{
		GID:gID,
		MinX:minX,
		MaxX:maxX,
		MinY:minY,
		MaxY:maxY,
		playerIDs:make(map[int]interface{}),
	}
}
//给格子添加一个玩家
func (g *Grid) Add(playerID int,player interface{}){
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	g.playerIDs[playerID]=player
}
//得到当前格子所有玩家的ID
func (g *Grid) GetPlayerIDs() (playerIDs []int){
	g.pIDLock.Lock()
	g.pIDLock.Unlock()
	for playerID,_:=range g.playerIDs{
		playerIDs=append(playerIDs,playerID)
	}
	return
}
//调试打印格子信息方法
func (g *Grid) String()string{
	return fmt.Sprintf("Grid id : %d,minX:%d,maxX:%d , minY:%d, maxY:%d, playerIDs:%v\n",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)

}
