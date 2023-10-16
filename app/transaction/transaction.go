package transaction

import (
	"context"
	"fmt"
	"time"
	"transactionapp/model"
	"transactionapp/store"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=$GOFILE -destination=transaction_mock.go -package=$GOPACKAGE

type App interface {
	InsertTransaction(ctx context.Context, transaction *model.Transaction) (string, error)
	InsertTransactions(ctx context.Context, transaction []*model.Transaction) ([]string, error)
	GetTransactions(ctx context.Context, transactionIds []string) ([]*model.TransactionResponse, error)
	GetTransactionsByPeriod(ctx context.Context, startDate, endDate string) ([]*model.TransactionResponse, error)
	GetTransactionsByPeriodEpoch(ctx context.Context, startDate, endDate int64) ([]*model.TransactionResponse, error)
}

type appImpl struct {
	stores *store.Container
	log    logrus.Logger
}

func NewAppTransaction(stores *store.Container, log logrus.Logger) App {
	return &appImpl{
		stores: stores,
		log:    log,
	}
}

func (a appImpl) InsertTransaction(ctx context.Context, transaction *model.Transaction) (string, error) {
	return a.stores.Transaction.InsertTransaction(ctx, transaction)
}

func convertDateToString(t time.Time) string {
	switch t.Month() {
	case time.January:
		return fmt.Sprintf("%d-01-%d", t.Year(), t.Day())
	case time.February:
		return fmt.Sprintf("%d-02-%d", t.Year(), t.Day())
	case time.March:
		return fmt.Sprintf("%d-03-%d", t.Year(), t.Day())
	case time.April:
		return fmt.Sprintf("%d-04-%d", t.Year(), t.Day())
	case time.May:
		return fmt.Sprintf("%d-05-%d", t.Year(), t.Day())
	case time.June:
		return fmt.Sprintf("%d-06-%d", t.Year(), t.Day())
	case time.July:
		return fmt.Sprintf("%d-07-%d", t.Year(), t.Day())
	case time.August:
		return fmt.Sprintf("%d-08-%d", t.Year(), t.Day())
	case time.September:
		return fmt.Sprintf("%d-09-%d", t.Year(), t.Day())
	case time.October:
		return fmt.Sprintf("%d-10-%d", t.Year(), t.Day())
	case time.November:
		return fmt.Sprintf("%d-11-%d", t.Year(), t.Day())
	default:
		return fmt.Sprintf("%d-12-%d", t.Year(), t.Day())
	}
}

func (a appImpl) InsertTransactions(ctx context.Context, transaction []*model.Transaction) ([]string, error) {
	for _, t := range transaction {
		if len(t.PurchaseDate) == 0 {
			currentTime := time.Now()
			t.CreatedAt = currentTime.Unix()
			t.PurchaseDate = convertDateToString(currentTime)
		} else {
			date, _ := formatDate(t.PurchaseDate)
			t.CreatedAt = date
		}
	}
	return a.stores.Transaction.InsertTransactions(ctx, transaction)
}

func (a appImpl) GetTransactions(ctx context.Context, transactionIds []string) ([]*model.TransactionResponse, error) {
	if len(transactionIds) == 0 {
		return nil, fmt.Errorf("failed to get transactions: empty id list")
	}

	return a.stores.Transaction.GetTransactionByIds(ctx, transactionIds)
}

func formatDate(date string) (int64, error) {
	d := fmt.Sprintf("%sT00:00:00", date)
	ts, err := time.ParseInLocation("2006-01-02T15:04:05", d, time.UTC)
	if err != nil {
		return 0, err
	}

	return ts.Unix(), nil
}

func (a appImpl) GetTransactionsByPeriod(ctx context.Context, startDate, endDate string) ([]*model.TransactionResponse, error) {
	sDate, err := formatDate(startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to format startDate: %w", err)
	}

	eDate, err := formatDate(endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to format endDate: %w", err)
	}

	return a.stores.Transaction.GetTransactionByDate(ctx, sDate, eDate)
}

func (a appImpl) GetTransactionsByPeriodEpoch(ctx context.Context, startDate, endDate int64) ([]*model.TransactionResponse, error) {
	return a.stores.Transaction.GetTransactionByDate(ctx, startDate, endDate)
}
