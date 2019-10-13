package account

import (
	"github.com/sgdr/wallet-service/internal/currency"
)

type Account struct {
	ID       uint
	Currency currency.Currency
	Owner    string
	Balance  uint64
}
