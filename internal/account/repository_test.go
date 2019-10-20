package account

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/sgdr/wallet-service/internal/config"
	"github.com/sgdr/wallet-service/internal/currency"
	"github.com/sgdr/wallet-service/internal/datasource"
	"github.com/sgdr/wallet-service/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestAccountRepository_Store(t *testing.T) {
	ctx := context.Background()
	l := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	ctx = logger.ToContext(ctx, l)
	d, err := datasource.CreateDataSource(ctx, config.DB{Name: "wallet", User: "wallet_user",
		Password: "wallet_user_pass", Host: "localhost", Port: "5432"})
	assert.Nil(t, err)
	a := Account{Owner: "Alice", Balance: 3445, Currency: currency.Currency{ID: 1}}
	r := NewRepository(d)
	err = r.Store(ctx, &a)
	assert.Nil(t, err)
	assert.True(t, a.ID > 0)
	accounts, err := r.FindAll(ctx)
	assert.Nil(t, err)
	//assert.Equal(t, 1, len(accounts))
	fmt.Printf("%+v", accounts[0])
	fmt.Println()
	accounts2, err2 := r.GetForUpdate(ctx, []string{"Alice"}, "USD")
	assert.Nil(t, err2)
	fmt.Println("----accounts-----")
	fmt.Printf("%+v", accounts2)
	rows, err := r.UpdateBalance(ctx, "Alice", accounts2[0].Currency.ID, 300)
	assert.Nil(t, err)
	fmt.Printf("%+v", rows)
}
