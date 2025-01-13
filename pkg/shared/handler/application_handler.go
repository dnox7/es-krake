package handler

import "log"

type ApplicationHTTPHandler struct {
	BaseHTTPHandler
}

func NewApplicationHTTPHandler(logger *log.Logger) *ApplicationHTTPHandler {
	return &ApplicationHTTPHandler{BaseHTTPHandler: BaseHTTPHandler{}}
}
