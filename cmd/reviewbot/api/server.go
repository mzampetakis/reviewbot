package api

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"golang.org/x/exp/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reviewbot/internal/database"
	"reviewbot/internal/domain/orders"
	"strconv"
	"sync"
	"syscall"
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

// The Server is used as a container for the most important dependencies.
type Server struct {
	Router      *mux.Router
	UserService *orders.Service
	App         *Application
}

// NewServer returns a pointer to a new Server.
func NewServer(userService *orders.Service, config *Application) *Server {
	server := &Server{
		Router:      mux.NewRouter().StrictSlash(true),
		UserService: userService,
		App:         config,
	}
	return server
}

func (mySrv *Server) ServeHTTP() error {
	srv := &http.Server{
		Addr:         net.JoinHostPort(mySrv.App.Config.BaseURL, strconv.Itoa(mySrv.App.Config.HttpPort)),
		Handler:      mySrv.routes(),
		ErrorLog:     slog.NewLogLogger(mySrv.App.Logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}
	shutdownErrorChan := make(chan error)

	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan
		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()
		shutdownErrorChan <- srv.Shutdown(ctx)
	}()

	mySrv.App.Logger.Info("starting server", slog.Group("server", "addr", srv.Addr))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrorChan
	if err != nil {
		mySrv.App.Logger.Error(err.Error())
		return err
	}

	mySrv.App.Logger.Info("stopped server", slog.Group("server", "addr", srv.Addr))
	mySrv.App.wg.Wait()
	return nil
}
