package transaction

import (
	"context"

	"github.com/jcpribeiro/TransactionApp/model"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -source=$GOFILE -destination=transaction_mock.go -package=$GOPACKAGE

type Store interface {
	InsertTransaction(ctx context.Context, transaction *model.Transaction) (string, error)
	InsertTransactions(ctx context.Context, transaction []*model.Transaction) ([]string, error)
	GetTransactionById(ctx context.Context, id string) (*model.TransactionResponse, error)
	GetTransactionByIds(ctx context.Context, ids []string) ([]*model.TransactionResponse, error)
	GetTransactionByDate(ctx context.Context, startDate, endDate int64) ([]*model.TransactionResponse, error)
}

type storeImpl struct {
	mongodbConReader *mongo.Database
	mongodbConWriter *mongo.Database
	log              logrus.Logger
}

const (
	descriptionMaxLength = 50
)

func NewStoreTransaction(mongodbConReader, mongodbConWriter *mongo.Database, log logrus.Logger) Store {
	return &storeImpl{
		mongodbConReader: mongodbConReader,
		mongodbConWriter: mongodbConWriter,
		log:              log,
	}
}

func setDescriptionMaxLength(description string) string {
	if len(description) > descriptionMaxLength {
		return description[:descriptionMaxLength]
	}

	return description
}

// Insert a new transaction
func (s storeImpl) InsertTransaction(ctx context.Context, transaction *model.Transaction) (string, error) {
	transaction.Description = setDescriptionMaxLength(transaction.Description)
	insertedId, err := s.mongodbConWriter.Collection("transaction").InsertOne(ctx, transaction)
	if err != nil {
		return "", err
	}

	return insertedId.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Insert an array of new transactions
func (s storeImpl) InsertTransactions(ctx context.Context, transaction []*model.Transaction) ([]string, error) {
	var exec []interface{}
	for _, t := range transaction {
		t.Description = setDescriptionMaxLength(t.Description)
		exec = append(exec, t)
	}

	insertedId, err := s.mongodbConWriter.Collection("transaction").InsertMany(ctx, exec)
	if err != nil {
		return []string{}, err
	}

	ids := []string{}
	for _, id := range insertedId.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID).Hex())
	}

	return ids, nil
}

// Get a single transaction info, filtering by id
func (s storeImpl) GetTransactionById(ctx context.Context, id string) (*model.TransactionResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := primitive.M{
		"_id": objectID,
	}

	var transaction *model.TransactionResponse
	err = s.mongodbConReader.Collection("transaction").FindOne(ctx, filter).Decode(&transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// Get a multiple transactions info, filtering by the ids array
func (s storeImpl) GetTransactionByIds(ctx context.Context, ids []string) ([]*model.TransactionResponse, error) {
	arrObjectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		arrObjectIDs = append(arrObjectIDs, objectID)
	}

	filter := primitive.M{
		"_id": primitive.M{"$in": arrObjectIDs},
	}

	transaction := make([]*model.TransactionResponse, 0, len(ids))
	cursor, err := s.mongodbConReader.Collection("transaction").Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// Get a multiple transactions info, filtering by date
func (s storeImpl) GetTransactionByDate(ctx context.Context, startDate, endDate int64) ([]*model.TransactionResponse, error) {
	filter := primitive.M{
		"created_at": primitive.M{"$gte": startDate, "$lt": endDate},
	}

	var transaction []*model.TransactionResponse
	cursor, err := s.mongodbConReader.Collection("transaction").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
