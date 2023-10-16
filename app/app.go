package app

import (
	"github.com/jcpribeiro/TransactionApp/app/fiscaldata"
	"github.com/jcpribeiro/TransactionApp/app/transaction"
	"github.com/jcpribeiro/TransactionApp/store"

	"github.com/sirupsen/logrus"
)

// Model container for exporting instantiated services
type Container struct {
	FiscalData  fiscaldata.App
	Transaction transaction.App
}

type Options struct {
	Log    logrus.Logger
	URL    string
	Stores *store.Container
}

// New creates a new instance of the services
func NewApp(opts Options) *Container {
	opts.Log.Info("Registered APP")

	return &Container{
		FiscalData:  fiscaldata.NewAppFiscalData(opts.URL, opts.Log),
		Transaction: transaction.NewAppTransaction(opts.Stores, opts.Log),
	}
}
