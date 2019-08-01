package main

import (
	"context"
	"fmt"
	"github.com/nettyrnp/url-shortener/config"
	"github.com/nettyrnp/url-shortener/http_service"
	"github.com/nettyrnp/url-shortener/log"
	"github.com/nettyrnp/url-shortener/router"
	"github.com/nettyrnp/url-shortener/storage_service"
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
		return nil, errors.Wrap(err, "creating storage http_service")
	}

	if a.httpService, err = http_service.NewHTTPService(ctx, a.config.HTTP); err != nil {
		return nil, errors.Wrap(err, "creating http http_service")
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	http.Handle("/", &router.RootHandler{Filename: "http_service/templates/index.html"})
	//http.Handle("/", &router.RootHandler{Filename: "index.html"})
	http.Handle("/v1/", &router.SvcHandler{})
	http.Handle("/app", a.httpService)

	addr := fmt.Sprintf("%v:%v", a.config.HTTP.Host, a.config.HTTP.Port)
	a.logger.Infof("Listening on %v\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		a.logger.Fatal("ListenAndServe:", err)
		return err
	}
	return nil
}
