package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/faizisyellow/indocoffee/internal/auth"
	"github.com/faizisyellow/indocoffee/internal/logger"
	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/faizisyellow/indocoffee/internal/uploader"
	"go.uber.org/zap"
)

type DBConf struct {
	Addr            string
	MaxOpenConn     int
	MaxIdleConn     int
	MaxIdleLifeTime string
	MaxLifeTime     string
}

type JwtConfig struct {
	SecretKey string
	Iss       string
	Sub       string

	// Unix time
	Exp int64
}

type Application struct {
	Port           string
	Host           string
	Env            string
	SwaggerUrl     string
	DbConfig       DBConf
	Services       service.Service
	JwtAuth        JwtConfig
	Authentication auth.Authenticator
	Upload         uploader.Uploader
	Logger         *zap.SugaredLogger
}

func (app *Application) Run(mux http.Handler) error {

	srv := http.Server{
		Addr:         net.JoinHostPort(app.Host, app.Port),
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {

		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		inf := strings.Builder{}
		inf.WriteString("signal cought, ")
		inf.WriteString(s.String())

		logger.Logger.Info(inf.String())

		shutdown <- srv.Shutdown(ctx)

	}()

	logger.Logger.Infow("server has started at", zap.String("host", app.Host), zap.String("port", app.Port), zap.String("env", app.Env))
	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	logger.Logger.Info("server has stopped at", zap.String("host", app.Host), zap.String("port", app.Port))

	return nil
}
