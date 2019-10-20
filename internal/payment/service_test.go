package payment

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/sgdr/wallet-service/internal/account"
	"github.com/sgdr/wallet-service/internal/config"
	"github.com/sgdr/wallet-service/internal/datasource"
	"github.com/sgdr/wallet-service/internal/db"
	"github.com/sgdr/wallet-service/internal/logger"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
)

// test needs in database "payment_test_db"
func TestPaymentService_CreatePayment(t *testing.T) {
	// setup
	ctx := context.Background()
	l := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	ctx = logger.ToContext(ctx, l)
	dTest, err := datasource.New(ctx, config.DB{Name: "payment_test_db", User: "wallet_user",
		Password: "wallet_user_pass", Host: "localhost", Port: "5432"})
	if err != nil {
		l.Log("error", err)
		t.FailNow()
	}

	err = db.CreateDbStructure(ctx, dTest, "../db/update_db_structure.sql")
	if err != nil {
		l.Log("error", err)
		t.FailNow()
	}
	defer dTest.Close()
	defer func() {
		dTest.Exec("drop table payments; drop table accounts; drop table currencies;")
	}()

	sqlQueries, err := ioutil.ReadFile("./testdata/data.sql")
	if err != nil {
		l.Log("error", err)
		t.FailNow()
	}
	tx, err := dTest.BeginTx(ctx, nil)
	if err != nil {
		l.Log("error", err)
		t.FailNow()
	}
	_, err = tx.ExecContext(ctx, string(sqlQueries))
	if err := tx.Commit(); err != nil {
		l.Log("error", err)
		t.FailNow()
	}
	accRep := account.NewRepository(dTest)
	paymentRep := NewRepository(dTest)
	ps := NewService(paymentRep, accRep)

	// Alise has 2000 USD and 3000 EUR
	// Bob has 300 USD
	// test case 1: try to transfer money from Alice to Bob twice
	// result: only one payment will be created
	var waitgroup = sync.WaitGroup{}
	waitgroup.Add(2)
	var err1 error
	var err2 error
	go func() {
		err1 = ps.CreatePayment(ctx, "Alice", "Bob", 1500, "USD")
		waitgroup.Done()
	}()
	go func() {
		err2 = ps.CreatePayment(ctx, "Alice", "Bob", 1500, "USD")
		waitgroup.Done()
	}()
	waitgroup.Wait()

	if err1 == nil {
		assert.NotNil(t, err2)
		assert.True(t, strings.Contains(err2.Error(), "unsufficient funds"))
	} else {
		assert.Nil(t, err2)
		assert.True(t, strings.Contains(err1.Error(), "unsufficient funds"))
	}

	payments, err := ps.GetAll(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(payments))
	assert.Equal(t, "Alice", payments[0].AccountFrom.Owner)
	assert.Equal(t, "Bob", payments[0].AccountTo.Owner)
	assert.Equal(t, "USD", payments[0].Currency.Ticker)
	assert.Equal(t, int64(1500), payments[0].Amount)
}
