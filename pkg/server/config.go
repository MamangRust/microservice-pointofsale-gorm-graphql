package server

import "time"

// Config defines the common configuration for all gRPC services
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	Port           int
	OtelEndpoint   string

	// DBCluster is the env prefix used to resolve this service's own
	// PostgreSQL connection (e.g. "DB_ROLE" -> DB_ROLE_HOST, DB_ROLE_NAME).
	// Empty falls back to the generic DB_* keys.
	DBCluster string

	// RedisCluster is the env prefix used to resolve this service's own Redis
	// connection. Empty falls back to the generic REDIS_* keys.
	RedisCluster string
}

// Default constants for gRPC server
const (
	DefaultMaxConcurrentConn = 1024
	DefaultWindowSize        = 16 * 1024 * 1024
	DefaultKeepaliveTime     = 20 * time.Second
	DefaultKeepaliveTimeout  = 5 * time.Second
	DefaultMinKeepaliveTime  = 5 * time.Second

	MonitoringInterval     = 30 * time.Second
	CleanupInterval        = 120 * time.Second
	CacheRefCountThreshold = 500

	ShutdownTimeout = 30 * time.Second

	RedisDialTimeout  = 5 * time.Second
	RedisReadTimeout  = 3 * time.Second
	RedisWriteTimeout = 3 * time.Second
	RedisPoolSize     = 10
	RedisMinIdleConns = 3
)
