package utils

import (
	"context"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type HttpClientNativeInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpClient struct {
	Client HttpClientNativeInterface
	Logger *logrus.Entry
}

type HttpClientInterface interface {
	Do(ctx context.Context, req *http.Request, cfg HttpClientConfig) (*http.Response, error)
}

type HttpClientConfig struct {
	RetryTimes uint
}

func NewHttpClient(logger *logrus.Entry) HttpClientInterface {
	return &HttpClient{
		Client: &http.Client{},
		Logger: logger,
	}
}

func (h HttpClient) Do(
	ctx context.Context,
	req *http.Request,
	cfg HttpClientConfig,
) (*http.Response, error) {

}

func (h HttpClient) shouldRetry(err error, res *http.Response) bool {
	if res == nil ||
		res.StatusCode == http.StatusBadGateway ||
		res.StatusCode == http.StatusServiceUnavailable ||
		res.StatusCode == http.StatusRequestTimeout ||
		res.StatusCode == http.StatusGatewayTimeout ||
		os.IsTimeout(err) {

		return true
	}

	return false
}

func (h HttpClient) backoff(retries int) time.Duration {
	return time.Duration(math.Pow(2, float64(retries))) * time.Second
}
