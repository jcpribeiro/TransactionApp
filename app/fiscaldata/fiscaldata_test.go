package fiscaldata

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRatesOfExchange(t *testing.T) {
	t.Run("This test simulates the exchange rate search", func(t *testing.T) {
		expected := Data{
			CurrencyDescription: "Canada-Dollar",
			ExchangeRate:        1.34,
			RecordDate:          "2023-10-15",
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			result, _ := json.Marshal(RateOfExchangeResponse{
				Data: []Data{
					0: expected,
				},
			})
			w.Write(result)
		}))
		defer server.Close()

		testFiscalData := appImpl{
			url:    server.URL,
			client: server.Client(),
		}

		data, err := testFiscalData.GetRatesOfExchange("Canada-Dollar", "2023-10-15")

		assert.Equal(t, data, &expected)
		assert.NoError(t, err)
	})

	t.Run("This test simulates an error when searching for the exchange rate", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(400)
			result, _ := json.Marshal(RateOfExchangeResponse{})
			w.Write(result)
		}))

		defer server.Close()

		testFiscalData := appImpl{
			url:    server.URL,
			client: server.Client(),
		}

		data, err := testFiscalData.GetRatesOfExchange("Canada-Dollar", "2023-10-15")

		assert.Nil(t, data)
		assert.Error(t, err)
	})
}
