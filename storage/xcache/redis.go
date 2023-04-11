package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

type (
	RedisClient interface {
		Get(key string) (string, error)
		Set(key, value string) (int, error)
		SetNx(key, value string) (int, error)
		SetExpired(key, value string, expired int) (int, error)
		SetNxExpired(key, value string, expired int) (int, error)
		Del(key string) (int, error)
	}

	DefaultRedisClient struct {
		pool *redis.Pool
	}

	RedisConfig struct {
		host            string
		port            int
		Password        string
		maxIdle         int
		maxActive       int
		maxConnLifetime time.Duration
		idleTimeout     time.Duration
	}

	RedisOption interface {
		apply(o *RedisConfig)
	}

	hostOption        string
	portOption        int
	passwordOption    string
	maxIdleOption     int
	maxActiveOption   int
	maxLifetimeOption time.Duration
	idleTimeoutOption time.Duration
)

const (
	defaultHost            = "127.0.0.1"
	defaultPort            = 6379
	defaultMaxIdle         = 50
	defaultMaxActive       = 100
	defaultMaxConnLifetime = 3600 * time.Second
	defaultIdleTimeout     = 240 * time.Second
)

func NewRedisClient(options ...RedisOption) DefaultRedisClient {
	conf := &RedisConfig{
		host:            defaultHost,
		port:            defaultPort,
		maxIdle:         defaultMaxIdle,
		maxActive:       defaultMaxActive,
		maxConnLifetime: defaultMaxConnLifetime,
		idleTimeout:     defaultIdleTimeout,
	}
	for _, o := range options {
		o.apply(conf)
	}
	return DefaultRedisClient{pool: newPool(conf)}
}

func newPool(conf *RedisConfig) *redis.Pool {
	return &redis.Pool{
		MaxIdle:         conf.maxIdle,
		MaxActive:       conf.maxActive,
		MaxConnLifetime: conf.maxConnLifetime,
		IdleTimeout:     conf.idleTimeout,
		Dial:            func() (redis.Conn, error) { return redis.Dial("tcp", fmt.Sprintf("%s:%d", conf.host, conf.port)) },
		Wait:            true,
	}
}

func (h hostOption) apply(o *RedisConfig) { o.host = string(h) }

func (p portOption) apply(o *RedisConfig) { o.port = int(p) }

func (p passwordOption) apply(o *RedisConfig) { o.Password = string(p) }

func (p maxIdleOption) apply(o *RedisConfig) { o.maxIdle = int(p) }

func (p maxActiveOption) apply(o *RedisConfig) { o.maxActive = int(p) }

func (p maxLifetimeOption) apply(o *RedisConfig) { o.maxConnLifetime = time.Duration(p) }

func (p idleTimeoutOption) apply(o *RedisConfig) { o.idleTimeout = time.Duration(p) }

func SetHost(h string) RedisOption { return hostOption(h) }

func SetPort(port int) RedisOption { return portOption(port) }

func SetPassword(pwd string) RedisOption { return passwordOption(pwd) }

func SetMaxIdle(idle int) RedisOption { return maxIdleOption(idle) }

func SetMaxActive(active int) RedisOption { return maxActiveOption(active) }

func SetMaxLifetime(l time.Duration) RedisOption { return maxLifetimeOption(l) }

func SetIdleTimeout(t time.Duration) RedisOption { return idleTimeoutOption(t) }

func (r DefaultRedisClient) Get(key string) (string, error) {
	return r.DoString("get")
}

func (r DefaultRedisClient) Set(key, value string) (int, error) {
	return r.DoInt("set", key, value)
}

func (r DefaultRedisClient) SetNx(key, value string) (int, error) {
	return r.DoInt("setnx", key, value)
}

func (r DefaultRedisClient) SetExpired(key, value string, expired int) (int, error) {
	return r.DoInt("set", key, value, "ex", expired)
}

func (r DefaultRedisClient) SetNxExpired(key, value string, expired int) (int, error) {
	return r.DoInt("setnx", key, value, "ex", expired)
}

func (r DefaultRedisClient) Del(key string) (int, error) {
	return r.DoInt("del", key)
}

func (r DefaultRedisClient) DoString(command string, args ...interface{}) (string, error) {
	conn := r.getConn()
	defer r.closeConn(conn)
	return redis.String(conn.Do(command, args...))
}

func (r DefaultRedisClient) DoInt(command string, args ...interface{}) (int, error) {
	conn := r.getConn()
	defer r.closeConn(conn)
	return redis.Int(conn.Do(command, args...))
}

func (r DefaultRedisClient) getConn() redis.Conn {
	return r.pool.Get()
}

func (r DefaultRedisClient) closeConn(conn redis.Conn) {
	_ = conn.Close()
}
