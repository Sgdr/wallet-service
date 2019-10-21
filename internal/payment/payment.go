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

type Direction string

const (
	Incoming Direction = "incoming"
	Outgoing Direction = "outgoing"
)

type ResponseItem struct {
	Amount    int64     `json:"amount"`
	Account   string    `json:"account"`
	Currency  string    `json:"currency"`
	Direction Direction `json:"direction"`
}

type CreatePaymentRequest struct {
	AccountFrom string `json:"accountFrom"`
	AccountTo   string `json:"accountTo"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
}

type ListRequest struct {
	ForAccount string `json:"forAccount"`
}

func (p *Payment) ToResponseItem(accountRequestFor string) ResponseItem {
	if p.AccountFrom.Owner == accountRequestFor {
		return ResponseItem{Account: p.AccountTo.Owner, Amount: p.Amount, Currency: p.Currency.Ticker, Direction: Outgoing}
	}
	return ResponseItem{Account: p.AccountFrom.Owner, Amount: p.Amount, Currency: p.Currency.Ticker, Direction: Incoming}
}
