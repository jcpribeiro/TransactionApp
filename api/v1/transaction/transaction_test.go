package transaction

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"transactionapp/app"
	"transactionapp/app/fiscaldata"
	"transactionapp/app/transaction"
	"transactionapp/internal/cache"
	"transactionapp/internal/validate"
	"transactionapp/model"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type strucTest struct {
	echo           *echo.Echo
	transactionApp *transaction.MockApp
	fiscalDataApp  *fiscaldata.MockApp
	cache          *cache.MockCache
}

func setUpTest(t *testing.T) strucTest {
	ctrl := gomock.NewController(t)
	echo := echo.New()
	echo.Validator = validate.New()
	transactionApp := transaction.NewMockApp(ctrl)
	fiscalDataApp := fiscaldata.NewMockApp(ctrl)
	cache := cache.NewMockCache(ctrl)

	return strucTest{
		echo:           echo,
		transactionApp: transactionApp,
		fiscalDataApp:  fiscalDataApp,
		cache:          cache,
	}
}

func TestInsertTransactions(t *testing.T) {
	t.Run("this test simulate a successful transaction insert", func(t *testing.T) {
		testObj := setUpTest(t)
		value := []*model.Transaction{
			0: {
				PurchaseAmount: 23.7,
				Description:    "Test",
				PurchaseDate:   "2023-10-15",
			},
		}
		body, _ := json.Marshal(value)
		req := httptest.NewRequest(http.MethodPost, "/v1/transaction", strings.NewReader(string(body)))
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		idList := []string{"652d34910a8fc425116b84d9"}
		testObj.transactionApp.EXPECT().InsertTransactions(gomock.Any(), gomock.Any()).Return(idList, nil)

		h := handler{
			apps: &app.Container{
				FiscalData:  testObj.fiscalDataApp,
				Transaction: testObj.transactionApp,
			},
		}

		ctx := testObj.echo.NewContext(req, rec)
		err := h.insertTransactions(ctx)

		var resp interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp)
		assert.NoError(t, err)
	})

	t.Run("this test simulate an error during a transaction insert", func(t *testing.T) {
		testObj := setUpTest(t)
		value := []*model.Transaction{
			0: {
				PurchaseAmount: 23.7,
				Description:    "Test",
				PurchaseDate:   "2023-10-15",
			},
		}
		body, _ := json.Marshal(value)
		req := httptest.NewRequest(http.MethodPost, "/v1/transaction", strings.NewReader(string(body)))
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		idList := []string{}
		testObj.transactionApp.EXPECT().InsertTransactions(gomock.Any(), gomock.Any()).Return(idList, errors.New("an error has ocurred"))

		h := handler{
			apps: &app.Container{
				FiscalData:  testObj.fiscalDataApp,
				Transaction: testObj.transactionApp,
			},
		}

		ctx := testObj.echo.NewContext(req, rec)
		err := h.insertTransactions(ctx)

		var resp interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Empty(t, resp)
		assert.Error(t, err)
	})
}

func TestGetTransactions(t *testing.T) {
	t.Run("This test simulates the process for obtaining transaction information", func(t *testing.T) {
		testObj := setUpTest(t)
		payload := []*model.TransactionResponse{
			0: {
				PurchaseAmount: 23.70,
				Description:    "Test1",
				PurchaseDate:   "2023-10-15",
			},
			1: {
				PurchaseAmount: 25.00,
				Description:    "Test2",
				PurchaseDate:   "2023-10-14",
			},
		}
		req := httptest.NewRequest(http.MethodGet, "/v1/transaction", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		testObj.transactionApp.EXPECT().GetTransactions(gomock.Any(), gomock.Any()).Return(payload, nil)
		testObj.cache.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any())
		testObj.fiscalDataApp.EXPECT().GetRatesOfExchange(gomock.Any(), gomock.Any()).Return(&fiscaldata.Data{
			CurrencyDescription: "Canada-Dollar",
			ExchangeRate:        1.23,
			RecordDate:          "2023-10-15",
		}, nil).AnyTimes()
		testObj.cache.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		h := handler{
			apps: &app.Container{
				FiscalData:  testObj.fiscalDataApp,
				Transaction: testObj.transactionApp,
			},
			cache: testObj.cache,
		}

		ctx := testObj.echo.NewContext(req, rec)
		ctx.QueryParams().Add("ids", "652d34910a8fc425116b84d9")
		ctx.QueryParams().Add("currency", "Canada-Dollar")
		err := h.getTransactions(ctx)

		var resp interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp)
		assert.NoError(t, err)
	})

	t.Run("This test simulates an error when obtaining transaction information", func(t *testing.T) {
		testObj := setUpTest(t)
		req := httptest.NewRequest(http.MethodGet, "/v1/transaction", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		testObj.cache.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any())
		testObj.transactionApp.EXPECT().GetTransactions(gomock.Any(), gomock.Any()).Return(nil, errors.New("an error has ocurred"))

		h := handler{
			apps: &app.Container{
				FiscalData:  testObj.fiscalDataApp,
				Transaction: testObj.transactionApp,
			},
			cache: testObj.cache,
		}

		ctx := testObj.echo.NewContext(req, rec)
		ctx.QueryParams().Add("ids", "652d34910a8fc425116b84d9")
		ctx.QueryParams().Add("currency", "Canada-Dollar")
		err := h.getTransactions(ctx)

		var resp interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Empty(t, resp)
		assert.Error(t, err)
	})
}

func TestGetTransactionsByPeriod(t *testing.T) {
	t.Run("This test simulates the process for obtaining transaction information by date", func(t *testing.T) {
		testObj := setUpTest(t)
		payload := []*model.TransactionResponse{
			0: {
				PurchaseAmount: 23.70,
				Description:    "Test1",
				PurchaseDate:   "2023-10-15",
			},
			1: {
				PurchaseAmount: 25.00,
				Description:    "Test2",
				PurchaseDate:   "2023-10-14",
			},
		}
		req := httptest.NewRequest(http.MethodGet, "/v1/transaction/period", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		testObj.transactionApp.EXPECT().GetTransactionsByPeriod(gomock.Any(), gomock.Any(), gomock.Any()).Return(payload, nil)
		testObj.fiscalDataApp.EXPECT().GetRatesOfExchange(gomock.Any(), gomock.Any()).Return(&fiscaldata.Data{
			CurrencyDescription: "Canada-Dollar",
			ExchangeRate:        1.23,
			RecordDate:          "2023-10-15",
		}, nil).Times(2)

		h := handler{
			apps: &app.Container{
				FiscalData:  testObj.fiscalDataApp,
				Transaction: testObj.transactionApp,
			},
			cache: testObj.cache,
		}

		ctx := testObj.echo.NewContext(req, rec)
		ctx.QueryParams().Add("ids", "652d34910a8fc425116b84d9")
		ctx.QueryParams().Add("currency", "Canada-Dollar")
		ctx.QueryParams().Add("startDate", "2023-10-12")
		ctx.QueryParams().Add("endDate", "2023-10-15")
		err := h.getTransactionsByPeriod(ctx)

		var resp interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp)
		assert.NoError(t, err)
	})

	t.Run("This test simulates an error when obtaining transaction information by date", func(t *testing.T) {
		testObj := setUpTest(t)
		req := httptest.NewRequest(http.MethodGet, "/v1/transaction/period", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		testObj.transactionApp.EXPECT().GetTransactionsByPeriod(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("an error has ocurred"))

		h := handler{
			apps: &app.Container{
				FiscalData:  testObj.fiscalDataApp,
				Transaction: testObj.transactionApp,
			},
			cache: testObj.cache,
		}

		ctx := testObj.echo.NewContext(req, rec)
		ctx.QueryParams().Add("ids", "652d34910a8fc425116b84d9")
		ctx.QueryParams().Add("currency", "Canada-Dollar")
		ctx.QueryParams().Add("startDate", "2023-10-12")
		ctx.QueryParams().Add("endDate", "2023-10-15")
		err := h.getTransactionsByPeriod(ctx)

		var resp interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Empty(t, resp)
		assert.Error(t, err)
	})
}

func TestGetTransactionsByPeriodEpoch(t *testing.T) {
	t.Run("This test simulates the process for obtaining transaction information by period", func(t *testing.T) {
		testObj := setUpTest(t)
		payload := []*model.TransactionResponse{
			0: {
				PurchaseAmount: 23.70,
				Description:    "Test1",
				PurchaseDate:   "2023-10-15",
			},
			1: {
				PurchaseAmount: 25.00,
				Description:    "Test2",
				PurchaseDate:   "2023-10-14",
			},
		}
		req := httptest.NewRequest(http.MethodGet, "/v1/transaction/epoch-period", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		testObj.transactionApp.EXPECT().GetTransactionsByPeriodEpoch(gomock.Any(), gomock.Any(), gomock.Any()).Return(payload, nil)
		testObj.fiscalDataApp.EXPECT().GetRatesOfExchange(gomock.Any(), gomock.Any()).Return(&fiscaldata.Data{
			CurrencyDescription: "Canada-Dollar",
			ExchangeRate:        1.23,
			RecordDate:          "2023-10-15",
		}, nil).Times(2)

		h := handler{
			apps: &app.Container{
				FiscalData:  testObj.fiscalDataApp,
				Transaction: testObj.transactionApp,
			},
			cache: testObj.cache,
		}

		ctx := testObj.echo.NewContext(req, rec)
		ctx.QueryParams().Add("ids", "652d34910a8fc425116b84d9")
		ctx.QueryParams().Add("currency", "Canada-Dollar")
		ctx.QueryParams().Add("startDate", "1697150153")
		ctx.QueryParams().Add("endDate", "1697409353")
		err := h.getTransactionsByPeriodEpoch(ctx)

		var resp interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp)
		assert.NoError(t, err)
	})

	t.Run("This test simulates an error when obtaining transaction information by date", func(t *testing.T) {
		testObj := setUpTest(t)
		req := httptest.NewRequest(http.MethodGet, "/v1/transaction/epoch-period", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		testObj.transactionApp.EXPECT().GetTransactionsByPeriodEpoch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("an error has ocurred"))

		h := handler{
			apps: &app.Container{
				FiscalData:  testObj.fiscalDataApp,
				Transaction: testObj.transactionApp,
			},
			cache: testObj.cache,
		}

		ctx := testObj.echo.NewContext(req, rec)
		ctx.QueryParams().Add("ids", "652d34910a8fc425116b84d9")
		ctx.QueryParams().Add("currency", "Canada-Dollar")
		ctx.QueryParams().Add("startDate", "1697150153")
		ctx.QueryParams().Add("endDate", "1697409353")
		err := h.getTransactionsByPeriodEpoch(ctx)

		var resp interface{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Empty(t, resp)
		assert.Error(t, err)
	})
}
