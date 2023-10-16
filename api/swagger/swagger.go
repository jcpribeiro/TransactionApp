package swagger

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	docs "transactionapp/docs"
)

// Options struct de opções para a criação de uma instancia do swagger
type Options struct {
	Group *echo.Group
}

// Register group item check
func Register(opts Options) {

	docs.SwaggerInfo.Title = "Swagger TransactionApp API"
	docs.SwaggerInfo.Description = "Swagger with TransactionApp API routes and parameter usage models"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	opts.Group.GET("", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/swagger/index.html")
	})

	opts.Group.GET("/*", func(c echo.Context) error {

		return echoSwagger.WrapHandler(c)
	})
}
