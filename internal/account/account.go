package account

import (
	"github.com/sgdr/wallet-service/internal/currency"
)

type Account struct {
	ID       uint
	Currency currency.Currency
	Owner    string
	Balance  int64
}

type ResponseItem struct {
	ID       string `json:"id"`
	Currency string `json:"currency"`
}

func (a *Account) ToResponseItem() ResponseItem {
	return ResponseItem{ID: a.Owner, Currency: a.Currency.Ticker}
}
