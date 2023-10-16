package api

import (
	v1 "github.com/jcpribeiro/TransactionApp/api/v1"
	"github.com/jcpribeiro/TransactionApp/app"
	"github.com/jcpribeiro/TransactionApp/internal/cache"

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
