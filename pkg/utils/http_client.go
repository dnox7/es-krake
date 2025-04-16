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

type (
	HttpClient interface {
		Do(ctx context.Context, req *http.Request, cfg HttpClientConfig) (*http.Response, error)
	}

	httpClient struct {
		client *http.Client
		logger *log.Logger
	}

	ReqOpt struct {
		CanLog                      bool
		CanLogRequestBody           bool
		CanLogResponseBody          bool
		CanLogRequestBodyOnlyError  bool
		CanLogResponseBodyOnlyError bool
		LoggedRequestBody           []string
		LoggedResponseBody          []string
		markedQueryParamKeys        []string
	}

	ClientOpt struct {
		RetryTimes            uint
		MaxIdleConnsPerHost   int
		Timeout               time.Duration
		ResponseHeaderTimeout time.Duration
		ServiceName           *string
	}
)

func NewHttpClient(serviceName string) HttpClient {
	return &httpClient{
		client: &http.Client{},
		logger: log.With("service", serviceName),
	}
}

func (h *httpClient) Do(
	ctx context.Context,
	req *http.Request,
	cfg HttpClientConfig,
) (*http.Response, error) {
	req = req.WithContext(ctx)
	logger := h.logger.
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

		res, err := h.client.Do(req)
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

func (h *httpClient) shouldRetry(err error, res *http.Response) bool {
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

func (h *httpClient) backoff(retries int) time.Duration {
	return time.Duration(math.Pow(2, float64(retries))) * time.Second
}

func (h *httpClient) outputLog(ctx context.Context, statusCode int, err error, fields map[string]interface{}) {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}

	if err != nil || statusCode/100 == 5 {
		h.logger.With(args...).Error(ctx, err.Error())
		return
	}

	msg := "Request completed with HttpClient"
	log.With(args...).Info(ctx, msg)
}
