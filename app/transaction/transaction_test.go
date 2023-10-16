package transaction

import (
	"context"
	"errors"
	"testing"
	"github.com/jcpribeiro/TransactionApp/model"
	"github.com/jcpribeiro/TransactionApp/store"
	"github.com/jcpribeiro/TransactionApp/store/transaction"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type structTest struct {
	storesMock *transaction.MockStore
	appTest    App
}

func setUptest(t *testing.T) structTest {
	ctrl := gomock.NewController(t)
	storesMock := transaction.NewMockStore(ctrl)

	return structTest{
		storesMock: storesMock,
		appTest: NewAppTransaction(&store.Container{
			Transaction: storesMock,
		}, *logrus.New()),
	}
}

func TestInsertTransaction(t *testing.T) {
	ctx := context.Background()

	t.Run("this test simulate a successful transaction insert", func(t *testing.T) {
		testObj := setUptest(t)
		expectedId := "652d34910a8fc425116b84d9"
		testObj.storesMock.EXPECT().InsertTransaction(ctx, &model.Transaction{
			PurchaseAmount: 23.70,
			Description:    "Test",
			PurchaseDate:   "2023-10-15",
		}).Return(expectedId, nil)

		id, err := testObj.appTest.InsertTransaction(ctx, &model.Transaction{
			PurchaseAmount: 23.70,
			Description:    "Test",
			PurchaseDate:   "2023-10-15",
		})

		assert.Equal(t, id, expectedId)
		assert.NoError(t, err)
	})

	t.Run("this test simulate an error during a transaction insert", func(t *testing.T) {
		testObj := setUptest(t)
		testObj.storesMock.EXPECT().InsertTransaction(ctx, &model.Transaction{
			PurchaseAmount: 23.70,
			Description:    "Test",
			PurchaseDate:   "2023-10-15",
		}).Return("", errors.New("an error has ocurred"))

		id, err := testObj.appTest.InsertTransaction(ctx, &model.Transaction{
			PurchaseAmount: 23.70,
			Description:    "Test",
			PurchaseDate:   "2023-10-15",
		})

		assert.Equal(t, id, "")
		assert.Error(t, err)
	})
}

func TestInsertTransactions(t *testing.T) {
	ctx := context.Background()

	t.Run("this test simulate a successful transactions insert", func(t *testing.T) {
		testObj := setUptest(t)
		expectedIds := []string{"652d34910a8fc425116b84d9", "652d34910a8fc425116b84d8"}
		payload := []*model.Transaction{
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
		testObj.storesMock.EXPECT().InsertTransactions(ctx, payload).Return(expectedIds, nil)

		ids, err := testObj.appTest.InsertTransactions(ctx, payload)

		assert.Equal(t, ids, expectedIds)
		assert.NoError(t, err)
	})

	t.Run("this test simulate a successful transactions insert", func(t *testing.T) {
		testObj := setUptest(t)
		expectedIds := []string{"652d34910a8fc425116b84d9", "652d34910a8fc425116b84d8"}
		payload := []*model.Transaction{
			0: {
				PurchaseAmount: 23.70,
				Description:    "Test1",
			},
			1: {
				PurchaseAmount: 25.00,
				Description:    "Test2",
			},
		}
		testObj.storesMock.EXPECT().InsertTransactions(ctx, payload).Return(expectedIds, nil)

		ids, err := testObj.appTest.InsertTransactions(ctx, payload)

		assert.Equal(t, ids, expectedIds)
		assert.NoError(t, err)
	})

	t.Run("this test simulate an error during a transactions insert", func(t *testing.T) {
		testObj := setUptest(t)
		payload := []*model.Transaction{
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
		testObj.storesMock.EXPECT().InsertTransactions(ctx, payload).Return([]string{}, errors.New("an error has ocurred"))

		id, err := testObj.appTest.InsertTransactions(ctx, payload)

		assert.Equal(t, id, []string{})
		assert.Error(t, err)
	})
}

func TestGetTransactions(t *testing.T) {
	ctx := context.Background()

	t.Run("This test simulates the process for obtaining transaction information", func(t *testing.T) {
		testObj := setUptest(t)
		idList := []string{"652d34910a8fc425116b84d9", "652d34910a8fc425116b84d8"}
		expectedResponse := []*model.TransactionResponse{
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
		testObj.storesMock.EXPECT().GetTransactionByIds(ctx, idList).Return(expectedResponse, nil)

		resp, err := testObj.appTest.GetTransactions(ctx, idList)

		assert.Equal(t, resp, expectedResponse)
		assert.NoError(t, err)
	})

	t.Run("This test simulates an error when obtaining transaction information", func(t *testing.T) {
		testObj := setUptest(t)
		idList := []string{"652d34910a8fc425116b84d9", "652d34910a8fc425116b84d8"}
		testObj.storesMock.EXPECT().GetTransactionByIds(ctx, idList).Return(nil, errors.New("an error has ocurred"))

		resp, err := testObj.appTest.GetTransactions(ctx, idList)

		assert.Nil(t, resp)
		assert.Error(t, err)
	})

	t.Run("This test simulates an error when obtaining transaction information", func(t *testing.T) {
		testObj := setUptest(t)
		idList := []string{}

		resp, err := testObj.appTest.GetTransactions(ctx, idList)

		assert.Nil(t, resp)
		assert.Error(t, err)
	})
}

func TestGetTransactionsByPeriod(t *testing.T) {
	ctx := context.Background()

	t.Run("This test simulates the process for obtaining transaction information by date", func(t *testing.T) {
		testObj := setUptest(t)
		startDate := "2023-10-12"
		endDate := "2023-10-15"
		sDate, _ := formatDate(startDate)
		eDate, _ := formatDate(endDate)
		expectedResponse := []*model.TransactionResponse{
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
		testObj.storesMock.EXPECT().GetTransactionByDate(ctx, sDate, eDate).Return(expectedResponse, nil)

		resp, err := testObj.appTest.GetTransactionsByPeriod(ctx, startDate, endDate)

		assert.Equal(t, resp, expectedResponse)
		assert.NoError(t, err)
	})

	t.Run("This test simulates an error when obtaining transaction information by date", func(t *testing.T) {
		testObj := setUptest(t)
		startDate := "2023-10-12"
		endDate := "2023-10-15"
		testObj.storesMock.EXPECT().GetTransactionByDate(ctx, gomock.Any(), gomock.Any()).Return(nil, errors.New("an error has ocurred"))

		resp, err := testObj.appTest.GetTransactionsByPeriod(ctx, startDate, endDate)

		assert.Nil(t, resp)
		assert.Error(t, err)
	})
}

func TestGetTransactionsByPeriodEpoch(t *testing.T) {
	ctx := context.Background()

	t.Run("This test simulates the process for obtaining transaction information by period", func(t *testing.T) {
		testObj := setUptest(t)
		var sDate int64 = 1697150153
		var eDate int64 = 1697409353
		expectedResponse := []*model.TransactionResponse{
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
		testObj.storesMock.EXPECT().GetTransactionByDate(ctx, sDate, eDate).Return(expectedResponse, nil)

		resp, err := testObj.appTest.GetTransactionsByPeriodEpoch(ctx, sDate, eDate)

		assert.Equal(t, resp, expectedResponse)
		assert.NoError(t, err)
	})

	t.Run("This test simulates an error when obtaining transaction information by period", func(t *testing.T) {
		testObj := setUptest(t)
		var sDate int64 = 1697150153
		var eDate int64 = 1697409353
		testObj.storesMock.EXPECT().GetTransactionByDate(ctx, sDate, eDate).Return(nil, errors.New("an error has ocurred"))

		resp, err := testObj.appTest.GetTransactionsByPeriodEpoch(ctx, sDate, eDate)

		assert.Nil(t, resp)
		assert.Error(t, err)
	})
}
