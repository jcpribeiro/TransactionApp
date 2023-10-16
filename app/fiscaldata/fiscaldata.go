package fiscaldata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=$GOFILE -destination=fiscaldata_mock.go -package=$GOPACKAGE

type App interface {
	GetRatesOfExchange(currencyDescription, transactionDate string) (*Data, error)
}

type appImpl struct {
	url    string
	client *http.Client
	log    logrus.Logger
}

type Data struct {
	CurrencyDescription string  `json:"country_currency_desc"`
	ExchangeRate        float64 `json:"exchange_rate,string"`
	RecordDate          string  `json:"record_date"`
}

type RateOfExchangeResponse struct {
	Data []Data `json:"data"`
}

func NewAppFiscalData(url string, log logrus.Logger) App {
	return &appImpl{
		url:    url,
		log:    log,
		client: newHttpClient(),
	}
}

func newHttpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}

func formatUrl(baseUrl, currencyDescription, startDate, endDate string) string {
	endpoint := "v1/accounting/od/rates_of_exchange"
	fields := "fields=country_currency_desc,exchange_rate,record_date"
	filter := fmt.Sprintf("filter=country_currency_desc:eq:%s,record_date:lte:%s,record_date:gte:%s", currencyDescription, endDate, startDate)
	sort := "sort=-record_date"
	formatedUrl := fmt.Sprintf("%s/%s?%s&%s&%s", baseUrl, endpoint, fields, filter, sort)
	return formatedUrl
}

func (e *RateOfExchangeResponse) getExchangeValue() *Data {
	if len(e.Data) > 0 {
		return &e.Data[0]
	}

	return nil
}

func convertDate(date string) (string, error) {
	layout := "2006-01-02T15:04:05"
	value := fmt.Sprintf("%sT00:00:00", date)

	ts, err := time.Parse(layout, value)
	if err != nil {
		return "", fmt.Errorf("failed to convert date: %w", err)
	}

	ts = ts.AddDate(0, -6, 0)
	date = fmt.Sprintf("%d-%d-%d", ts.Year(), ts.Month(), ts.Day())

	if ts.Month() < 10 {
		date = fmt.Sprintf("%d-0%d-%d", ts.Year(), ts.Month(), ts.Day())
	}

	return date, nil
}

func (a appImpl) GetRatesOfExchange(currencyDescription, transactionDate string) (*Data, error) {
	convertedDate, err := convertDate(transactionDate)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Get(formatUrl(a.url, currencyDescription, convertedDate, transactionDate))
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body information: %w", err)
	}

	responseData := &RateOfExchangeResponse{}
	err = json.Unmarshal(body, responseData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal body response: %w", err)
	}

	return responseData.getExchangeValue(), nil
}
