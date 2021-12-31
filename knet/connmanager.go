package knet

import (
	"errors"
	"kinx/kiface"
	"log"
	"sync"
)

type ConnectionManager struct {
	connections map[uint32]kiface.IConnection // map记录所有的连接
	mu          sync.RWMutex
}

func (cm *ConnectionManager) GetConns() map[uint32]kiface.IConnection {
	return cm.connections
}

// 添加连接
func (cm *ConnectionManager) Add(connection kiface.IConnection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.connections[connection.GetConnectionID()] = connection
}

// 删除连接
func (cm *ConnectionManager) Remove(connection kiface.IConnection) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, ok := cm.connections[connection.GetConnectionID()]; !ok {
		log.Println("connection", connection.GetConnectionID(), "do not exits!")
		return errors.New("remove nil pointer connection")
	}

	delete(cm.connections, connection.GetConnectionID())
	return nil
}

// 根据connID获取一个连接
func (cm *ConnectionManager) Get(id uint32) (kiface.IConnection, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if _, ok := cm.connections[id]; !ok {
		log.Println("connection", id, "do not exit!")
		return nil, errors.New("get nil pointer connection")
	}
	return cm.connections[id], nil
}

// 删除所有连接
func (cm *ConnectionManager) Clear() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for id, conn := range cm.connections {
		// 停止conn
		conn.StopWithNotConnMgr()

		// 删除conn
		delete(cm.connections, id)
	}
}

// 获取当前连接总个数
func (cm *ConnectionManager) Len() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return len(cm.connections)
}

func NewConnMgr() kiface.IConnectionManager {
	return &ConnectionManager{
		connections: make(map[uint32]kiface.IConnection),
		mu:          sync.RWMutex{},
	}
}
