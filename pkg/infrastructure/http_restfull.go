package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pech/es-krake/pkg/shared/utils"
	"strings"

	"github.com/sirupsen/logrus"
)

type RestfullHttpClient interface {
	Get(ctx context.Context, req Request) (statusCode int, res []byte, err error)
	Post(ctx context.Context, req Request) (statusCode int, res []byte, err error)
}

type Request struct {
	Endpoint string
	Data     map[string]interface{}
	Query    string
	Headers  map[string]string
}

type restfullHttpClient struct {
	baseUrl string
	client  http.Client
	logger  *logrus.Entry
}

type RestfullHttpClientConfig struct {
	TagName string
	BaseUrl string
}

func NewRestfullHttpClient(conf RestfullHttpClientConfig) RestfullHttpClient {
	logger := NewLogger()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return restfullHttpClient{
		baseUrl: conf.BaseUrl,
		logger:  logrus.WithField("service", conf.TagName+"-api"),
		client: http.Client{
			Transport: utils.NewLoggedHttpTransport(
				logger.WithField("service", conf.TagName+"-api"),
				utils.HttpTransportExternalType,
			),
		},
	}
}

func (c restfullHttpClient) Get(ctx context.Context, req Request) (int, []byte, error) {
	return c.do(ctx, Do{
		Method:   http.MethodGet,
		Endpoint: req.Endpoint,
		Query:    req.Query,
		Body:     strings.NewReader(""),
		Headers:  req.Headers,
	})
}

func (c restfullHttpClient) Post(ctx context.Context, req Request) (int, []byte, error) {
	var body *strings.Reader

	switch req.Headers["Content-Type"] {
	case "application/json":
		bodyStr, err := json.Marshal(req.Data)
		if err != nil {
			return 0, nil, err
		}
		body = strings.NewReader(string(bodyStr))

	case "application/x-www-form-urlencoded":
		values := url.Values{}
		for k, v := range req.Data {
			values.Set(k, fmt.Sprintf("%+v", v))
		}
		body = strings.NewReader(values.Encode())
	}

	return c.do(ctx, Do{
		Method:   http.MethodPost,
		Endpoint: req.Endpoint,
		Query:    req.Query,
		Body:     body,
		Headers:  req.Headers,
	})
}

type Do struct {
	Method   string
	Endpoint string
	Query    string
	Body     *strings.Reader
	Headers  map[string]string
}

func (c restfullHttpClient) do(ctx context.Context, req Do) (int, []byte, error) {
	url := c.baseUrl + req.Endpoint
	if req.Query != "" {
		url += "?" + req.Query
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, req.Body)
	if err != nil {
		return 0, nil, err
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	res, err := c.client.Do(httpReq)
	if err != nil {
		return 0, nil, err
	}

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, nil, err
	}

	return res.StatusCode, bodyBytes, nil
}
