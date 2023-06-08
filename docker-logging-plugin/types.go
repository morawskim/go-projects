package main

import "github.com/docker/docker/daemon/logger"

type CapabilitiesResponse struct {
	Err string
	Cap logger.Capability
}

type StartLoggingRequest struct {
	File string
	Info logger.Info
}

type StopLoggingRequest struct {
	File string
}

type response struct {
	Err string
}
