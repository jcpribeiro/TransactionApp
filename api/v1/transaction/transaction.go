package transaction

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"transactionapp/app"
	"transactionapp/internal/cache"
	"transactionapp/internal/util"
	"transactionapp/model"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Register group transaction
func Register(g *echo.Group, apps *app.Container, cache cache.Cache) {
	h := &handler{
		apps:  apps,
		cache: cache,
	}

	g.POST("", h.insertTransactions)
	g.GET("", h.getTransactions)
	g.GET("/period", h.getTransactionsByPeriod)
	g.GET("/epoch-period", h.getTransactionsByPeriodEpoch)
}

type handler struct {
	apps  *app.Container
	cache cache.Cache
}

// insertTransactions swagger document
// @Summary Store a purchase transaction
// @Tags transaction
// @Accept  json
// @Produce  json
// @Param transaction body []model.Transaction true "add new transaction"
// @Success 200 {array} string
// @Failure 401 {object} string
// @Router /v1/transaction [post]
func (h *handler) insertTransactions(c echo.Context) error {
	var transactions []*model.Transaction
	if err := c.Bind(&transactions); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid message",
		})
	}

	response, err := h.apps.Transaction.InsertTransactions(c.Request().Context(), transactions)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string][]string{
		"ids": response,
	})
}

// getTransactions swagger document
// @Summary Retrive stored a purchase transaction
// @Tags transaction
// @Accept  json
// @Produce  json
// @Param ids query string true "Transactions ids. If more than one id is provided, it must be separated by a comma. E.g. id1,id2"
// @Param currency query string true "Currency ids. E.g. Argentina-Peso"
// @Success 200 {array} model.TransactionResponse
// @Failure 401 {object} string
// @Router /v1/transaction [get]
func (h *handler) getTransactions(c echo.Context) error {
	params := new(model.GetTransactionParams)

	if err := c.Bind(params); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid query params",
		})
	}

	if err := c.Validate(params); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing url params",
		})
	}

	ids := strings.Split(strings.ReplaceAll(params.Ids, " ", ""), ",")
	var err error
	response := make([]*model.TransactionResponse, 0, len(ids))
	for _, id := range ids {
		var value *model.TransactionResponse
		h.cache.Get(c.Request().Context(), fmt.Sprintf("%s:transaction:%s", id, params.Currency), &value)
		if value != nil {
			response = append(response, value)
		}
	}

	if len(response) != len(ids) {
		response, err = h.apps.Transaction.GetTransactions(c.Request().Context(), ids)
		if err != nil {
			return err
		}

		for _, r := range response {
			data, err := h.apps.FiscalData.GetRatesOfExchange(params.Currency, r.PurchaseDate)
			if err != nil {
				return err
			}
			r.PurchaseAmount = util.RoundFloat(r.PurchaseAmount, 2)
			r.ExchangeRate = util.RoundFloat(data.ExchangeRate, 2)
			r.ConvertedPurchaseAmount = util.RoundFloat(r.PurchaseAmount*data.ExchangeRate, 2)
			h.cache.Set(c.Request().Context(), fmt.Sprintf("%s:transaction:%s", r.Id, params.Currency), r, 5*time.Minute)
		}
	}

	return c.JSON(http.StatusOK, map[string][]*model.TransactionResponse{
		"ids": response,
	})
}

// getTransactionsByPeriod swagger document
// @Summary Retrive stored a purchase transaction by period
// @Tags transaction
// @Accept  json
// @Produce  json
// @Param startDate query string true "Period start date. E.g. 2023-10-12"
// @Param endDate query string true "Period end date. E.g. 2023-10-14"
// @Param currency query string true "Currency ids. E.g. Argentina-Peso"
// @Success 200 {array} model.TransactionResponse
// @Failure 401 {object} string
// @Router /v1/transaction/period [get]
func (h *handler) getTransactionsByPeriod(c echo.Context) error {
	params := new(model.GetTransactionParamsByPeriod)

	if err := c.Bind(params); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid query params",
		})
	}

	if err := c.Validate(params); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing url params",
		})
	}

	response, err := h.apps.Transaction.GetTransactionsByPeriod(c.Request().Context(), params.StartDate, params.EndDate)
	if err != nil {
		return err
	}

	for _, r := range response {
		data, err := h.apps.FiscalData.GetRatesOfExchange(params.Currency, r.PurchaseDate)
		if err != nil {
			return err
		}
		r.PurchaseAmount = util.RoundFloat(r.PurchaseAmount, 2)
		r.ExchangeRate = util.RoundFloat(data.ExchangeRate, 2)
		r.ConvertedPurchaseAmount = util.RoundFloat(r.PurchaseAmount*data.ExchangeRate, 2)
	}

	return c.JSON(http.StatusOK, map[string][]*model.TransactionResponse{
		"ids": response,
	})
}

// getTransactionsByPeriodEpoch swagger document
// @Summary Retrive stored a purchase transaction by period using epoch format
// @Tags transaction
// @Accept  json
// @Produce  json
// @Param startDate query string true "Period start date. E.g. 1697150153"
// @Param endDate query string true "Period end date. E.g. 1697409353"
// @Param currency query string true "Currency ids. E.g. Argentina-Peso"
// @Success 200 {array} model.TransactionResponse
// @Failure 401 {object} string
// @Router /v1/transaction/epoch-period [get]
func (h *handler) getTransactionsByPeriodEpoch(c echo.Context) error {
	params := new(model.GetTransactionParamsByPeriodEpoch)

	if err := c.Bind(params); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid query params",
		})
	}

	if err := c.Validate(params); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing url params",
		})
	}

	response, err := h.apps.Transaction.GetTransactionsByPeriodEpoch(c.Request().Context(), params.StartDate, params.EndDate)
	if err != nil {
		return err
	}

	for _, r := range response {
		data, err := h.apps.FiscalData.GetRatesOfExchange(params.Currency, r.PurchaseDate)
		if err != nil {
			return err
		}
		r.PurchaseAmount = util.RoundFloat(r.PurchaseAmount, 2)
		r.ExchangeRate = util.RoundFloat(data.ExchangeRate, 2)
		r.ConvertedPurchaseAmount = util.RoundFloat(r.PurchaseAmount*data.ExchangeRate, 2)
	}

	return c.JSON(http.StatusOK, map[string][]*model.TransactionResponse{
		"ids": response,
	})
}
