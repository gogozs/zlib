package xsql

import (
	"fmt"
	"github.com/gogozs/zlib/xlog"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	defaultMaxIdle         = 4
	defaultMaxActive       = 8
	defaultMaxConnLifetime = 3600 * time.Second
	defaultIdleTimeout     = 240 * time.Second
)

type (
	SQLConfig struct {
		Host     string
		Username string
		Password string
		Port     int
		Dbname   string
	}
	option struct {
		MaxIdle         int
		MaxOpen         int
		MaxConnLifetime time.Duration
		MaxIdleTime     time.Duration
	}
	OptionFunc func(o *option)
)

func WithMaxIdle(maxIdle int) OptionFunc {
	return func(o *option) {
		o.MaxIdle = maxIdle
	}
}

func WithMaxActive(maxActive int) OptionFunc {
	return func(o *option) {
		o.MaxIdle = maxActive
	}
}

func WithMaxConnLifetime(maxConnLifetime time.Duration) OptionFunc {
	return func(o *option) {
		o.MaxConnLifetime = maxConnLifetime
	}
}

func WithIdleTimeout(idleTimeout time.Duration) OptionFunc {
	return func(o *option) {
		o.MaxIdleTime = idleTimeout
	}
}

func NewDB(sqlConfig *SQLConfig, opts ...OptionFunc) (DB, error) {
	xlog.Info("config: %s", sqlConfig)
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s",
		sqlConfig.Username, sqlConfig.Password, sqlConfig.Host, sqlConfig.Port, sqlConfig.Dbname)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	sqlOptions := &option{
		MaxIdle:         defaultMaxIdle,
		MaxOpen:         defaultMaxActive,
		MaxConnLifetime: defaultMaxConnLifetime,
		MaxIdleTime:     defaultIdleTimeout,
	}
	for _, f := range opts {
		f(sqlOptions)
	}
	db.SetMaxIdleConns(sqlOptions.MaxIdle)
	db.SetMaxOpenConns(sqlOptions.MaxOpen)
	db.SetConnMaxLifetime(sqlOptions.MaxConnLifetime)
	db.SetConnMaxIdleTime(sqlOptions.MaxIdleTime)
	return db, nil
}

func NewLogDB(sqlConfig *SQLConfig, options ...OptionFunc) (DB, error) {
	db, err := NewDB(sqlConfig, options...)
	if err != nil {
		return nil, err
	}
	return WrapLog(db), nil
}
