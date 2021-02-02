package webreal

import "time"

// 客户端配置
type Config struct {
	WriteWait       time.Duration // 写等待时间
	WriteBufferSize int           // 写缓冲长度
	PongWait        time.Duration // 心跳等待时间
	PingInterval    time.Duration // 心跳频率
	MaxMessageSize  int64         // 最大消息字节数
}

func DefaultConfig() *Config {
	return &Config{
		WriteWait:       10 * time.Second,
		WriteBufferSize: 1024,
		PongWait:        60 * time.Second,
		PingInterval:    54 * time.Second,
		MaxMessageSize:  524288, // 512KB
	}
}
