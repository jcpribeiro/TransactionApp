package store

import (
	"github.com/jcpribeiro/TransactionApp/store/transaction"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Container struct {
	Transaction transaction.Store
}

type Options struct {
	MongodbConReader *mongo.Database
	MongodbConWriter *mongo.Database
	Log              logrus.Logger
}

func NewStore(opts Options) *Container {
	logrus.Info("Registered STORE")
	return &Container{
		Transaction: transaction.NewStoreTransaction(opts.MongodbConReader, opts.MongodbConWriter, opts.Log),
	}
}
