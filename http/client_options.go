package http

import (
	"github.com/go-resty/resty/v2"
	"time"
)

type (
	ClientOptions struct {
		proxy          *Proxy
		maxRetries     int
		waitTime       time.Duration
		maxWaitTime    time.Duration
		retryAfterFunc resty.RetryAfterFunc
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

func WithRetryAfterFunc(retryAfterFunc resty.RetryAfterFunc) ClientOptionFunc {
	return func(o *ClientOptions) {
		o.retryAfterFunc = retryAfterFunc
	}
}
