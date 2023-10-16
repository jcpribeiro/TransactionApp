package mongodb

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//go:generate mockgen -source=$GOFILE -destination=mongodb_mock.go -package=$GOPACKAGE

type MongoDB interface {
	Connect() *mongo.Database
}

type mongodb struct {
	url      string
	database string
	isReader bool
}

func NewMongoDB(url, database string, isReader bool) MongoDB {
	return &mongodb{
		url:      url,
		database: database,
		isReader: isReader,
	}
}

func (impl *mongodb) Connect() *mongo.Database {
	opts := options.Client()
	opts.ApplyURI(impl.url)

	if impl.isReader {
		opts.SetReadPreference(readpref.SecondaryPreferred())
	}

	c, err := mongo.NewClient(opts)
	if err != nil {
		logrus.Fatal("error on connect to mongodb:  ", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = c.Connect(ctx)
	if err != nil {
		logrus.Fatal("error on connect to mongodb: ", err.Error())
	}

	err = c.Ping(ctx, readpref.Primary())
	if err != nil {
		logrus.Fatal("error on connect to mongodb: ", err.Error())
	}

	return c.Database(impl.database)
}
