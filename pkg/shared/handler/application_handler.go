package handler

import "github.com/sirupsen/logrus"

type ApplicationHTTPHandler struct {
	BaseHTTPHandler
}

func NewApplicationHTTPHandler(logger *logrus.Logger) *ApplicationHTTPHandler {
	return &ApplicationHTTPHandler{BaseHTTPHandler: BaseHTTPHandler{Logger: logger}}
}
