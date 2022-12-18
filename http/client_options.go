package http

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type (
	ClientOptions struct {
		proxy              *Proxy
		maxRetries         int
		waitTime           time.Duration
		maxWaitTime        time.Duration
		retryConditionFunc resty.RetryConditionFunc
	}
	Proxy struct {
		Host     string
		Port     int
		Username string
		Password string
	}
	ClientOptionFunc func(o *ClientOptions)
)

func WithProxy(proxy *Proxy) ClientOptionFunc {
	return func(o *ClientOptions) {
		o.proxy = proxy
	}
}

func WithMaxRetries(maxRetries int) ClientOptionFunc {
	return func(o *ClientOptions) {
		o.maxRetries = maxRetries
	}
}

func WithWaitTime(waitTime time.Duration) ClientOptionFunc {
	return func(o *ClientOptions) {
		o.waitTime = waitTime
	}
}
func WithMaxWaitTime(maxWaitTime time.Duration) ClientOptionFunc {
	return func(o *ClientOptions) {
		o.maxWaitTime = maxWaitTime
	}
}
