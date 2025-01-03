package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
	req = req.WithContext(ctx)
	logger := h.Logger.
		WithContext(req.Context()).
		WithFields(logrus.Fields{
			"url":      req.URL.String(),
			"method":   req.Method,
			"protocol": req.Proto,
		})

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
			logger.WithError(err).Error("could not get body of response")
			return nil, err
		}

		bodyRes, err := io.ReadAll(res.Body)
		if err != nil {
			logger.WithError(err).Error("could not read body of response")
			return nil, err
		}

		res.Body = io.NopCloser(bytes.NewBuffer(bodyRes))
		logger = logger.WithFields(logrus.Fields{
			"duration":  fmt.Sprintf("%vms", time.Since(before).Milliseconds()),
			"retry_num": retries,
		})

		if res.StatusCode/100 == 5 {
			logger.WithError(err).Error("http call failed")
		} else {
			logger.WithField("status", res.Status).Info("http call successfully")
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
