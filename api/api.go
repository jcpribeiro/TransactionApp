package api

import (
	v1 "transactionapp/api/v1"
	"transactionapp/app"
	"transactionapp/internal/cache"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Options struct for creating an instance of the routes
type Options struct {
	Group *echo.Group
	Apps  *app.Container
	Cache cache.Cache
}

// Register api instance
func Register(opts Options) {
	v1.Register(opts.Group, opts.Apps, opts.Cache)
	// healthz.Register(opts.Root, opts.Apps)

	logrus.Info("Registered API")
}
