package transaction

import (
	"context"
	"testing"
	"time"
	"transactionapp/model"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

type structTest struct {
	mt *mtest.T
}

func prepareTest(t *testing.T) *structTest {
	mt := mtest.New(t, mtest.NewOptions().DatabaseName("test").ClientType(mtest.Mock))
	return &structTest{
		mt: mt,
	}
}

func TestInsertTransaction(t *testing.T) {
	testObj := prepareTest(t)
	ctx := context.Background()

	testObj.mt.Run("this test simulate a successful transaction insert", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateSuccessResponse())
		storeTest := NewStoreTransaction(t.DB, t.DB, *logrus.New())

		id, err := storeTest.InsertTransaction(ctx, &model.Transaction{
			PurchaseAmount: 23.70,
			Description:    "Test - 12345678910111213141516171819202122324252627282930313233",
			PurchaseDate:   "2023-10-15",
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	testObj.mt.Run("this test simulate an error during a transaction insert", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code: 2,
		}))
		storeTest := NewStoreTransaction(t.DB, t.DB, *logrus.New())

		id, err := storeTest.InsertTransaction(ctx, &model.Transaction{
			PurchaseAmount: 23.70,
			Description:    "Test",
			PurchaseDate:   "2023-10-15",
		})

		assert.Error(t, err)
		assert.Empty(t, id)
	})
}

func TestInsertTransactions(t *testing.T) {
	testObj := prepareTest(t)
	ctx := context.Background()

	testObj.mt.Run("this test simulate a successful transactions insert", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateSuccessResponse())
		storeTest := NewStoreTransaction(t.DB, t.DB, *logrus.New())

		ids, err := storeTest.InsertTransactions(ctx, []*model.Transaction{
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
		})

		assert.NoError(t, err)
		assert.Equal(t, len(ids), 2)
	})

	testObj.mt.Run("this test simulate an error during a transactions insert", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code: 2,
		}))
		storeTest := NewStoreTransaction(t.DB, t.DB, *logrus.New())

		id, err := storeTest.InsertTransactions(ctx, []*model.Transaction{
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
		})

		assert.Error(t, err)
		assert.Empty(t, id)
	})
}

func TestGetTransactionById(t *testing.T) {
	testObj := prepareTest(t)
	ctx := context.Background()

	testObj.mt.Run("This test simulates the process for obtaining transaction information", func(t *mtest.T) {
		id := primitive.NewObjectID()
		ct := time.Now()
		expected := model.TransactionResponse{
			Id:             id.Hex(),
			PurchaseAmount: 25.00,
			Description:    "Test",
			CreatedAt:      ct.Unix(),
		}
		t.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expected.Id},
			{Key: "purchase_amount", Value: expected.PurchaseAmount},
			{Key: "description", Value: expected.Description},
			{Key: "created_at", Value: expected.CreatedAt},
		}))

		storeTest := NewStoreTransaction(t.DB, t.DB, *logrus.New())

		transactionTest, err := storeTest.GetTransactionById(ctx, id.Hex())

		assert.NoError(t, err)
		assert.Equal(t, transactionTest, &expected)
	})

	testObj.mt.Run("This test simulates an error when obtaining transaction information", func(t *mtest.T) {
		id := primitive.NewObjectID()
		t.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code: 2,
		}))

		storeTest := NewStoreTransaction(t.DB, t.DB, *logrus.New())

		transactionTest, err := storeTest.GetTransactionById(ctx, id.Hex())

		assert.Error(t, err)
		assert.Nil(t, transactionTest)
	})

	testObj.mt.Run("This test simulates an error when obtaining transaction information - invalid hex id", func(t *mtest.T) {
		storeTest := NewStoreTransaction(t.DB, t.DB, *logrus.New())

		transactionTest, err := storeTest.GetTransactionById(ctx, "test")

		assert.Error(t, err)
		assert.Nil(t, transactionTest)
	})
}

func TestGetTransactionByIds(t *testing.T) {
	testObj := prepareTest(t)
	ctx := context.Background()

	testObj.mt.Run("This test simulates the process for obtaining transactions information", func(t *mtest.T) {
		id_1 := primitive.NewObjectID()
		id_2 := primitive.NewObjectID()
		idList := []string{id_1.Hex(), id_2.Hex()}
		ct := time.Now()
		expected := []*model.TransactionResponse{
			0: {
				Id:             id_1.Hex(),
				PurchaseAmount: 25.00,
				Description:    "Test_1",
				CreatedAt:      ct.Unix(),
			},
			1: {
				Id:             id_2.Hex(),
				PurchaseAmount: 30.00,
				Description:    "Test_2",
				CreatedAt:      ct.Unix(),
			},
		}
		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expected[0].Id},
			{Key: "purchase_amount", Value: expected[0].PurchaseAmount},
			{Key: "description", Value: expected[0].Description},
			{Key: "created_at", Value: expected[0].CreatedAt},
		})

		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{Key: "_id", Value: expected[1].Id},
			{Key: "purchase_amount", Value: expected[1].PurchaseAmount},
			{Key: "description", Value: expected[1].Description},
			{Key: "created_at", Value: expected[1].CreatedAt},
		})

		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		t.AddMockResponses(first, second, killCursors)

		storeTest := NewStoreTransaction(t.DB, t.DB, *logrus.New())

		transactionTest, err := storeTest.GetTransactionByIds(ctx, idList)

		assert.NoError(t, err)
		assert.Equal(t, transactionTest, expected)
	})

	testObj.mt.Run("This test simulates an error when obtaining transactions information", func(t *mtest.T) {
		id_1 := primitive.NewObjectID()
		id_2 := primitive.NewObjectID()
		idList := []string{id_1.Hex(), id_2.Hex()}
		t.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code: 2,
		}))

		storeTest := NewStoreTransaction(t.DB, t.DB, *logrus.New())

		transactionTest, err := storeTest.GetTransactionByIds(ctx, idList)

		assert.Error(t, err)
		assert.Nil(t, transactionTest)
	})

	testObj.mt.Run("This test simulates an error when obtaining transactions information - invalid hex id", func(t *mtest.T) {
		id_1 := primitive.NewObjectID()
		idList := []string{"test", id_1.Hex()}

		storeTest := NewStoreTransaction(t.DB, t.DB, *logrus.New())

		transactionTest, err := storeTest.GetTransactionByIds(ctx, idList)

		assert.Error(t, err)
		assert.Nil(t, transactionTest)
	})
}

func TestGetTransactionByDate(t *testing.T) {
	testObj := prepareTest(t)
	ctx := context.Background()

	testObj.mt.Run("This test simulates the process for obtaining transactions information by date", func(t *mtest.T) {
		expected := []*model.TransactionResponse{
			0: {
				Id:             primitive.NewObjectID().Hex(),
				PurchaseAmount: 25.00,
				Description:    "Test_1",
				CreatedAt:      time.Now().Unix(),
			},
			1: {
				Id:             primitive.NewObjectID().Hex(),
				PurchaseAmount: 30.00,
				Description:    "Test_2",
				CreatedAt:      time.Now().Unix(),
			},
		}
		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expected[0].Id},
			{Key: "purchase_amount", Value: expected[0].PurchaseAmount},
			{Key: "description", Value: expected[0].Description},
			{Key: "created_at", Value: expected[0].CreatedAt},
		})

		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{Key: "_id", Value: expected[1].Id},
			{Key: "purchase_amount", Value: expected[1].PurchaseAmount},
			{Key: "description", Value: expected[1].Description},
			{Key: "created_at", Value: expected[1].CreatedAt},
		})

		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		t.AddMockResponses(first, second, killCursors)

		storeTest := NewStoreTransaction(t.DB, t.DB, *logrus.New())

		transactionTest, err := storeTest.GetTransactionByDate(ctx, 1697150153, 1697409353)

		assert.NoError(t, err)
		assert.Equal(t, transactionTest, expected)
	})

	testObj.mt.Run("This test simulates an error when obtaining transactions information by date", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code: 2,
		}))

		storeTest := NewStoreTransaction(t.DB, t.DB, *logrus.New())

		transactionTest, err := storeTest.GetTransactionByDate(ctx, 1697150153, 1697409353)

		assert.Error(t, err)
		assert.Nil(t, transactionTest)
	})
}
