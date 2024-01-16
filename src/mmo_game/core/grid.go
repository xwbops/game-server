package core

import (
	"fmt"
	"sync"
)

/*
	一个地图中的格子类
*/
type Grid struct {
	GID       int          //格子id
	MinX      int          //格子左边界坐标
	MaxX      int          //格子右边界坐标
	MinY      int          //格子上边界坐标
	MaxY      int          //格子下边界坐标
	playerIDs map[int]bool //当前格子内所有对象id
	pIDLock   sync.RWMutex //playerIDs保护锁

}

//初始化一个格子
func NewGrid(gId, minX, maxX, minY, MaxY int) *Grid {
	return &Grid{
		GID:       0,
		MinX:      0,
		MaxX:      0,
		MinY:      0,
		MaxY:      0,
		playerIDs: make(map[int]bool),
	}
}

//向当前格子中添加一个玩家
func (g *Grid) AddPlayer(playerId int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	g.playerIDs[playerId] = true
}

//向当前格子中删除一个玩家
func (g *Grid) RemovePlayer(playerId int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	delete(g.playerIDs, playerId)
}

//得到当前格子所有玩家
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()
	for id := range g.playerIDs {
		playerIDs = append(playerIDs, id)
	}
	return
}

//打印信息方法
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX: %d, maxX: %d, minY: %d, maxY: %d, playerIDs: %v", g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
