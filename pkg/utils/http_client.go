package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"pech/es-krake/pkg/log"
	"time"
)

type HttpClientNativeInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpClient struct {
	Client HttpClientNativeInterface
	Logger *log.Logger
}

type HttpClientInterface interface {
	Do(ctx context.Context, req *http.Request, cfg HttpClientConfig) (*http.Response, error)
}

type HttpClientConfig struct {
	RetryTimes uint
}

func NewHttpClient(serviceName string) HttpClientInterface {
	return &HttpClient{
		Client: &http.Client{},
		Logger: log.With("service", serviceName),
	}
}

func (h HttpClient) Do(
	ctx context.Context,
	req *http.Request,
	cfg HttpClientConfig,
) (*http.Response, error) {
	req = req.WithContext(ctx)
	logger := h.Logger.
		With(
			"url", req.URL.String(),
			"method", req.Method,
			"protocol", req.Proto,
		)

	var (
		res     *http.Response
		err     error
		retries = 0
	)

	for h.shouldRetry(err, res) && retries <= int(cfg.RetryTimes) {
		time.Sleep(h.backoff(retries))
		before := time.Now()

		res, err := h.Client.Do(req)
		if err != nil {
			logger.Error(
				ctx, "could not get body of response",
				"error", err,
			)
			return nil, err
		}

		bodyRes, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Error(
				ctx, "could not read body of response",
				"error", err,
			)
			return nil, err
		}

		res.Body = io.NopCloser(bytes.NewBuffer(bodyRes))
		logger = logger.With(
			"duration", fmt.Sprintf("%vms", time.Since(before).Milliseconds()),
			"retry_num", retries,
		)

		if res.StatusCode/100 == 5 {
			logger.Error(
				ctx, "http call failed",
				"error", err,
			)
		} else {
			logger.With("status", res.Status).Info(ctx, "http call successfully")
		}

		retries++
	}

	return res, err
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
