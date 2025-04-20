package httpclient

import (
	"net/http"
	"time"
)

type (
	ClientOpt struct{}

	clientOptSetter func(*clientOpt)

	clientOpt struct {
		client                *http.Client
		maxIdleConnsPerHost   int
		timeout               time.Duration
		responseHeaderTimeout time.Duration
		serviceName           string
	}
)

func (ClientOpt) Client(cli *http.Client) clientOptSetter {
	return func(co *clientOpt) {
		co.client = cli
	}
}

func (ClientOpt) Timeout(timeout time.Duration) clientOptSetter {
	return func(co *clientOpt) {
		co.timeout = timeout
	}
}

func (ClientOpt) MaxIdleConnsPerHost(conns int) clientOptSetter {
	return func(co *clientOpt) {
		co.maxIdleConnsPerHost = conns
	}
}

func (ClientOpt) ResponseHeaderTimeout(timeout time.Duration) clientOptSetter {
	return func(co *clientOpt) {
		co.responseHeaderTimeout = timeout
	}
}

func (ClientOpt) ServiceName(name string) clientOptSetter {
	return func(co *clientOpt) {
		co.serviceName = name
	}
}
