package v1

import (
	"transactionapp/app"
	"transactionapp/internal/cache"

	"transactionapp/api/v1/transaction"

	"github.com/labstack/echo/v4"
)

// Registers v1 routes
func Register(g *echo.Group, apps *app.Container, cache cache.Cache) {
	v1 := g.Group("/v1")

	transaction.Register(v1.Group("/transaction"), apps, cache)
}
