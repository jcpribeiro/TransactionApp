package v1

import (
	"github.com/jcpribeiro/TransactionApp/app"
	"github.com/jcpribeiro/TransactionApp/internal/cache"

	"github.com/jcpribeiro/TransactionApp/api/v1/transaction"

	"github.com/labstack/echo/v4"
)

// Registers v1 routes
func Register(g *echo.Group, apps *app.Container, cache cache.Cache) {
	v1 := g.Group("/v1")

	transaction.Register(v1.Group("/transaction"), apps, cache)
}
