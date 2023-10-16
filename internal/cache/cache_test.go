package cache

import (
	"context"
	"errors"
	"testing"
	"github.com/jcpribeiro/TransactionApp/model"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type structTest struct {
	db      *redis.Client
	mock    redismock.ClientMock
	appTest Cache
}

func setUpTest() structTest {
	db, mock := redismock.NewClientMock()
	appTest := NewCache(db, *logrus.New())

	return structTest{
		db:      db,
		mock:    mock,
		appTest: appTest,
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	t.Run("This test simulates an error when retrieving data from the cache", func(t *testing.T) {
		testObj := setUpTest()

		testObj.mock.ExpectGet("test").RedisNil()

		var value *model.TransactionResponse
		testObj.appTest.Get(ctx, "test", &value)
	})

	t.Run("This test simulates an error when retrieving data from the cache", func(t *testing.T) {
		testObj := setUpTest()

		testObj.mock.ExpectGet("test").SetErr(errors.New("an erro has ocurred"))

		var value *model.TransactionResponse
		testObj.appTest.Get(ctx, "test", &value)
	})

	t.Run("This test simulates retrieving data from the cache", func(t *testing.T) {
		testObj := setUpTest()
		expectedValue := "{\"key\":\"test-value\"}"
		testObj.mock.ExpectGet("test").SetVal(expectedValue)

		var value interface{}
		testObj.appTest.Get(ctx, "test", &value)

		assert.NotNil(t, value, expectedValue)
	})
}

func TestSet(t *testing.T) {
	ctx := context.Background()
	t.Run("This test simulates an error during the data set in the cache", func(t *testing.T) {
		testObj := setUpTest()

		testObj.mock.ExpectSet("key", "value", 1).RedisNil()

		err := testObj.appTest.Set(ctx, "key", "value", 1)
		assert.Error(t, err)
	})

	t.Run("This test simulates an error during the data set in the cache", func(t *testing.T) {
		testObj := setUpTest()

		testObj.mock.ExpectSet("key", "value", 1).SetErr(errors.New("an erro has ocurred"))

		err := testObj.appTest.Set(ctx, "key", "value", 1)
		assert.Error(t, err)
	})

	t.Run("This test simulates the data set in the cache", func(t *testing.T) {
		testObj := setUpTest()

		value, _ := marshalBinary("value")
		testObj.mock.ExpectSet("key", value, 1).SetVal("OK")

		err := testObj.appTest.Set(ctx, "key", "value", 1)
		assert.NoError(t, err)
	})
}
