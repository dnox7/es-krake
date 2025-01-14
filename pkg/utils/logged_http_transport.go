package utils

import (
	"io"
	"net/http"
	"pech/es-krake/pkg/log"
	"time"
)

type LoggedHttpTransportType string

const (
	HttpTransportExternalType LoggedHttpTransportType = "EXTERNAL"
	HttpTransportInternalType LoggedHttpTransportType = "INTERNAL"
)

type LoggedHttpTransport struct {
	http.RoundTripper
	logger *log.Logger
	Type   LoggedHttpTransportType
}

func NewLoggedHttpTransport(
	logFieldKey string,
	logFieldVal string,
	transportType LoggedHttpTransportType,
) http.RoundTripper {
	return LoggedHttpTransport{
		http.DefaultTransport,
		log.With(logFieldKey, logFieldVal),
		transportType,
	}
}

func (t LoggedHttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	before := time.Now()
	logger := t.logger.WithContext(req.Context())
	res, err := t.RoundTripper.RoundTrip(req)

	loggerRes := logger.With(
		"url", req.URL.String(),
		"method", req.Method,
		"protocol", req.Proto,
	)
	if err != nil {
		loggerRes.Error(req.Context(), "could not get response", "error", err)
		return nil, err
	}

	after := time.Now()
	var body []byte
	if req.Body != nil {
		bodyReq, err := req.GetBody()
		if err != nil {
			loggerRes.Error(req.Context(), "could not get request body", "error", err)
			return nil, err
		}

		body, err = io.ReadAll(io.LimitReader(bodyReq, 1024))
		if err != nil {
			loggerRes.Error(req.Context(), "could not read body of request", "error", err)
			return nil, err
		}

		defer bodyReq.Close()
	}

	bodyRes, err := io.ReadAll(res.Body)
	if err != nil {
		loggerRes.Error(req.Context(), "could not get body of response", "error", err)
		return nil, err
	}

	loggerRes = loggerRes.With(
		"duration", after.Sub(before).String(),
		"body_request", string(body),
		"body_response", string(bodyRes),
		"request_params", req.URL.Query(),
	)

	if res.StatusCode/100 == 5 {
		loggerRes.With("rootcause", t.Type).Error(req.Context(), err.Error())
	} else {
		loggerRes.With("status", res.Status).Info(req.Context(), "")
	}

	return res, err
}
