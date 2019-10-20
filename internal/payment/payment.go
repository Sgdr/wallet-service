package payment

import (
	"time"

	"github.com/sgdr/wallet-service/internal/account"
	"github.com/sgdr/wallet-service/internal/currency"
)

type Payment struct {
	ID          uint
	AccountFrom account.Account
	AccountTo   account.Account
	Amount      int64
	Currency    currency.Currency
	Date        time.Time
}
