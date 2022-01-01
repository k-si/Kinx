package knet

import "time"

const (
	DefaultIPVersion         = "tcp4"
	DefaultHost              = "127.0.0.1"
	DefaultTcpPort           = 8088
	DefaultMaxConnSize       = 1
	DefaultMaxPackageSize    = 1024 * 1024 // 1mb
	DefaultWorkerPoolSize    = 1
	DefaultMaxWorkerTaskSize = 100
	DefaultHeartRateInSecond = 30 * time.Second
	DefaultHeartFreshLevel   = 5
)

type Config struct {
	IPVersion         string
	Host              string
	TcpPort           int
	MaxConnSize       int    // 最多可有多少连接
	MaxPackageSize    uint32 // 一次可发送包的最大size
	WorkerPoolSize    uint32 // 多少个worker
	MaxWorkerTaskSize uint32 // 每个worker最多承载的任务数量
	HeartRateInSecond time.Duration
	HeartFreshLevel   uint32
}

func DefaultConfig() Config {
	return Config{
		IPVersion:         DefaultIPVersion,
		Host:              DefaultHost,
		TcpPort:           DefaultTcpPort,
		MaxConnSize:       DefaultMaxConnSize,
		MaxPackageSize:    DefaultMaxPackageSize,
		WorkerPoolSize:    DefaultWorkerPoolSize,
		MaxWorkerTaskSize: DefaultMaxWorkerTaskSize,
		HeartRateInSecond: DefaultHeartRateInSecond,
		HeartFreshLevel:   DefaultHeartFreshLevel,
	}
}
