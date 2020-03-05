package nets

import (
	"errors"
	"fmt"
	"sync"
	"topaza/interfaces"
)

// 连接管理模块
type ConnManager struct {
	// 连接集合
	connections map[uint32] interfaces.IConnection

	// 保护连接集合读写锁
	connLock sync.RWMutex
}

// 创建连接
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]interfaces.IConnection),
	}
}

// 添加连接
func (cm *ConnManager) Add(conn interfaces.IConnection) {
	// 保护共享资源， 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将 conn 加入到 ConnManager 中
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connection:", conn.GetConnID(), " add to connManager success, conn count:", cm.Len())
}

// 删除连接
func (cm *ConnManager) Remove(conn interfaces.IConnection) {
	// 保护共享资源， 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除连接
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connection:", conn.GetConnID(), " remove from connManager success, conn count:", cm.Len())
}

// 根据 ConnID 获取连接
func (cm *ConnManager) Get(connID uint32) (interfaces.IConnection, error) {
	// 保护共享资源， 加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")
}

// 得到当前连接总数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// 清除并终止所有连接
func (cm *ConnManager) ClearConn() {
	// 保护共享资源， 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除 conn 并停止 conn 工作
	for connID, conn := range cm.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.connections, connID)
	}

	fmt.Println("clear all connections success, conn count:", cm.Len())
}