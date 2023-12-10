package api

import (
	"golang.org/x/exp/slog"
	"reviewbot/internal/database"
	"sync"
	"time"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
	handlerDefaultTimeout = 20 * time.Second
)

type Application struct {
	Config ApplicationConfig
	DB     *database.DB
	Logger *slog.Logger
	wg     sync.WaitGroup
}

type ApplicationConfig struct {
	BaseURL  string
	HttpPort int
	DB       struct {
		DSN         string
		Automigrate bool
	}
}
