package webreal

import "time"

// 客户端配置
type Config struct {
	WriteWait       time.Duration // 写等待时间
	WriteChanBuffer int           // 写缓冲长度
	PongWait        time.Duration // 心跳等待时间
	PingPeriod      time.Duration // 心跳频率
	MaxMessageSize  int64         // 最大消息字节数
}

func DefaultConfig() *Config {
	return &Config{
		WriteWait:       10 * time.Second,
		WriteChanBuffer: 256,
		PongWait:        60 * time.Second,
		PingPeriod:      54 * time.Second,
		MaxMessageSize:  524288, // 512KB
	}
}
