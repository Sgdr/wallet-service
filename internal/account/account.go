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

type ResponseItem struct {
	ID       string `json:"id"`
	Balance  uint64 `json:"balance"`
	Currency string `json:"currency"`
}

func (a *Account) ToResponseItem() ResponseItem {
	return ResponseItem{ID: a.Owner, Balance: a.Balance, Currency: a.Currency.Ticker}
}
