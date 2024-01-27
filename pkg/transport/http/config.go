package http

import "time"

type Config struct {
	Host         string
	Port         uint16
	Silent       bool
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}
