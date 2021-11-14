package kiface

type IConnectionManager interface {
	// 添加连接
	Add(connection IConnection)

	// 删除连接
	Remove(connection IConnection) error

	// 根据connID获取一个连接
	Get(uint32) (IConnection, error)

	// 删除所有连接
	Clear()

	// 获取当前连接总个数
	Len() int
}
