package main

import (
	"context"
	"fmt"
	"github.com/nettyrnp/url-shortener/config"
	"github.com/nettyrnp/url-shortener/http"
	"github.com/nettyrnp/url-shortener/log"
	"github.com/nettyrnp/url-shortener/router"
	"github.com/nettyrnp/url-shortener/storage"
	"github.com/nettyrnp/url-shortener/util"
	"go.uber.org/zap"
	"net/http"

	"github.com/pkg/errors"
)

type App struct {
	config         config.Config
	logger         *zap.SugaredLogger
	httpService    *http_service.HTTPService
	storageService *storage_service.StorageService
}

var (
	logger = log.GetLogger()
)

func NewApp(ctx context.Context) (*App, error) {
	logger := log.GetLogger()
	defer logger.Sync() // flushes buffer, if any

	conf, err := config.GetConfig()
	util.Die(err)

	a := &App{
		config: *conf,
		logger: logger,
	}
	if a.storageService, err = storage_service.NewDBStorageService(ctx, a.config.Storage); err != nil {
		return nil, errors.Wrap(err, "creating storage http")
	}

	if a.httpService, err = http_service.NewHTTPService(ctx, a.config.HTTP); err != nil {
		return nil, errors.Wrap(err, "creating http http")
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	http.Handle("/", &router.RootHandler{Filename: "http/templates/index.html"})
	http.Handle("/v1/", &router.SvcHandler{})
	http.Handle("/app", a.httpService)

	addr := fmt.Sprintf("%v:%v", a.config.HTTP.Host, a.config.HTTP.Port)
	a.logger.Infof("Http Service is listening on %v", addr)
	return http.ListenAndServe(addr, nil)
}
