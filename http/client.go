package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gogozs/zlib/log"
)

type (
	Client struct {
		client *resty.Client
	}
)

var defaultClientOptions = ClientOptions{
	maxRetries:  3,
	waitTime:    time.Second * 5,
	maxWaitTime: time.Second * 20,
	retryConditionFunc: func(response *resty.Response, err error) bool {
		if response == nil ||
			response.StatusCode() == 0 ||
			(response.StatusCode() >= http.StatusLocked && response.StatusCode() < http.StatusNotExtended) {
			return true
		}
		return false
	},
}

var defaultRequestOptions = RequestOptions{
	headers: map[string]string{},
}

func NewHttpClient(options ...ClientOptionFunc) *Client {
	clientOptions := defaultClientOptions
	for _, o := range options {
		o(&clientOptions)
	}

	// Create a Resty Client
	client := resty.New()
	// Retries are configured per client
	client.
		// Set retry count to non zero to enable retries
		SetRetryCount(clientOptions.maxRetries).
		// You can override initial retry wait time.
		// Default is 100 milliseconds.
		SetRetryWaitTime(clientOptions.waitTime).
		// MaxWaitTime can be overridden as well.
		// Default is 2 seconds.
		SetRetryMaxWaitTime(clientOptions.maxWaitTime).
		AddRetryCondition(clientOptions.retryConditionFunc)

	c := &Client{client: client}
	c.SetProxy(clientOptions.proxy)

	return c
}

func (c *Client) SetProxy(proxy *Proxy) {
	if proxy == nil {
		return
	}
	proxyURL := fmt.Sprintf("http://%s:%s@%s:%d", proxy.Username, proxy.Password, proxy.Host, proxy.Port)
	log.InfoContext(context.Background(), "proxyURL: %s", proxyURL)
	c.client.SetProxy(proxyURL)
}

func (c *Client) Get(url string, opts ...RequestOptionFunc) (rsp *resty.Response, err error) {
	options := c.getRequestOptions(opts...)
	log.DebugContext(context.Background(), "HttpClient Start Get url: %s", url)
	defer func() {
		log.DebugContext(context.Background(), "HttpClient Start Get url done. err: %v", err)
	}()
	resp, err := c.client.R().
		EnableTrace().
		SetHeaders(options.headers).
		Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Post(url string, body interface{}, opts ...RequestOptionFunc) (rsp *resty.Response, err error) {
	options := c.getRequestOptions(opts...)
	log.DebugContext(context.Background(), "HttpClient Start Post url: %s", url)
	defer func() {
		log.DebugContext(context.Background(), "HttpClient Start Post url done. err: %v", err)
	}()
	resp, err := c.client.R().
		SetBody(body).
		EnableTrace().
		SetHeaders(options.headers).
		Post(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) getRequestOptions(opts ...RequestOptionFunc) RequestOptions {
	requestOptions := defaultRequestOptions
	for _, o := range opts {
		o(&requestOptions)
	}
	return requestOptions
}
