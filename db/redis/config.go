package redis

type Config struct {
	Host     string
	Port     int
	Password string
	Db       int
	// 是否启动redis
	Enabled  uint
	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int `mapstructure:"pool_size"`
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int `mapstructure:"minIdle_Conns"`
}
