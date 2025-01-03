package utils

import (
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type LoggedHttpTransportType string

const (
	HttpTransportExternalType LoggedHttpTransportType = "EXTERNAL"
	HttpTransportInternalType LoggedHttpTransportType = "INTERNAL"
)

type LoggedHttpTransport struct {
	http.RoundTripper
	logger *logrus.Entry
	Type   LoggedHttpTransportType
}

func NewLoggedHttpTransport(
	logger *logrus.Entry,
	transportType LoggedHttpTransportType,
) http.RoundTripper {
	return LoggedHttpTransport{
		http.DefaultTransport,
		logger,
		transportType,
	}
}

func (t LoggedHttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	before := time.Now()
	logger := t.logger.WithContext(req.Context())
	res, err := t.RoundTripper.RoundTrip(req)

	loggerRes := logger.WithFields(logrus.Fields{
		"url":      req.URL.String(),
		"method":   req.Method,
		"protocol": req.Proto,
	})
	if err != nil {
		loggerRes.Errorf("could not get response: %v", err)
		return nil, err
	}

	after := time.Now()
	var body []byte
	if req.Body != nil {
		bodyReq, err := req.GetBody()
		if err != nil {
			loggerRes.Errorf("could not get request body: %v", err)
			return nil, err
		}

		body, err = io.ReadAll(io.LimitReader(bodyReq, 1024))
		if err != nil {
			loggerRes.Errorf("could not read body of request: %v", err)
			return nil, err
		}

		defer bodyReq.Close()
	}

	bodyRes, err := io.ReadAll(res.Body)
	if err != nil {
		loggerRes.Errorf("could not get body of response: %v", err)
		return nil, err
	}

	loggerRes = loggerRes.WithFields(logrus.Fields{
		"duration":       after.Sub(before).String(),
		"body_request":   string(body),
		"body_response":  string(bodyRes),
		"request_params": req.URL.Query(),
	})

	if res.StatusCode/100 == 5 {
		loggerRes.WithField("rootcause", t.Type).Errorln(err)
	} else {
		loggerRes.WithField("status", res.Status).Info()
	}

	return res, err
}
