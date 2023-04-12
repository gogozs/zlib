package cache

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type (
	RedisClient interface {
		Get(key string) (string, error)
		Set(key, value string) (string, error)
		SetNx(key, value string) (int, error)
		SetExpired(key, value string, expired int) (string, error)
		SetNxExpired(key, value string, expired int) (string, error)
		Del(key string) (int, error)
		Incr(key string, value int) (int, error)
		Eval(script string, keys []string, argv []string) (interface{}, error)
		EvalSha(sha string, script string, keys []string, argv []string) (interface{}, error)
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

var _ RedisClient = (*DefaultRedisClient)(nil)

const (
	defaultHost            = "127.0.0.1"
	defaultPort            = 6379
	defaultMaxIdle         = 50
	defaultMaxActive       = 100
	defaultMaxConnLifetime = 3600 * time.Second
	defaultIdleTimeout     = 240 * time.Second
)

const noMatchingScriptErr = "NOSCRIPT No matching script. Please use EVAL."

func NewRedisClient(options ...RedisOption) RedisClient {
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
	return r.doString("get", key)
}

func (r DefaultRedisClient) Set(key, value string) (string, error) {
	return r.doString("set", key, value)
}

func (r DefaultRedisClient) SetNx(key, value string) (int, error) {
	return r.doInt("setnx", key, value)
}

func (r DefaultRedisClient) SetExpired(key, value string, expired int) (string, error) {
	return r.doString("set", key, value, "ex", expired)
}

func (r DefaultRedisClient) SetNxExpired(key, value string, expired int) (string, error) {
	return r.doString("set", key, value, "ex", expired, "nx")
}

func (r DefaultRedisClient) Del(key string) (int, error) {
	return r.doInt("del", key)
}

func (r DefaultRedisClient) Incr(key string, value int) (int, error) {
	conn := r.getConn()
	defer r.closeConn(conn)
	return redis.Int(conn.Do("INCRBY", key, value))
}

func (r DefaultRedisClient) Eval(script string, keys []string, argv []string) (interface{}, error) {
	conn := r.getConn()
	defer r.closeConn(conn)

	args := r.buildArgs(script, keys, argv)
	return conn.Do("EVAL", args...)
}

func (r DefaultRedisClient) EvalSha(sha string, script string, keys []string, argv []string) (interface{}, error) {
	conn := r.getConn()
	defer r.closeConn(conn)

	args := r.buildArgs(sha, keys, argv)
	res, err := conn.Do("EVALSHA", args...)
	if err != nil && err.Error() == noMatchingScriptErr {
		return r.Eval(script, keys, argv)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r DefaultRedisClient) buildArgs(target string, keys []string, argv []string) []interface{} {
	args := make([]any, 0, len(keys)+len(argv)+2)
	args = append(args, target, len(keys))
	for _, key := range keys {
		args = append(args, key)
	}
	for _, arg := range argv {
		args = append(args, arg)
	}
	return args
}

func (r DefaultRedisClient) doString(command string, args ...interface{}) (string, error) {
	conn := r.getConn()
	defer r.closeConn(conn)
	return redis.String(conn.Do(command, args...))
}

func (r DefaultRedisClient) doInt(command string, args ...interface{}) (int, error) {
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
