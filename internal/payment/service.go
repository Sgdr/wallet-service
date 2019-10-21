package payment

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/go-kit/kit/log/level"
	"github.com/sgdr/wallet-service/internal/account"
	"github.com/sgdr/wallet-service/internal/datasource"
	"github.com/sgdr/wallet-service/internal/logger"
)

type Service interface {
	GetAllForClient(ctx context.Context, client string) ([]Payment, error)
	CreatePayment(ctx context.Context, sender string, recipient string, amount int64, currency string) error
}

type paymentService struct {
	paymentRep Repository
	accountRep account.Repository
}

func NewService(paymentRep Repository, accountRep account.Repository) Service {
	return &paymentService{paymentRep: paymentRep, accountRep: accountRep}
}

func (p paymentService) GetAllForClient(ctx context.Context, client string) ([]Payment, error) {
	log := logger.FromContext(ctx)
	level.Info(log).Log("msg", "get all payments for client"+client)
	return p.paymentRep.FindAll(ctx, client)
}

func (p paymentService) CreatePayment(ctx context.Context, sender string, recipient string, amount int64, currency string) error {
	log := logger.FromContext(ctx)
	level.Info(log).Log("msg", fmt.Sprintf("create payment from %s to %s amount %d %s", sender, recipient, amount, currency))
	ds := datasource.GetDataSource()
	tx, err := ds.BeginTx(ctx, nil)
	if err != nil {
		level.Error(log).Log("msg", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			stack := string(debug.Stack())
			err := fmt.Errorf("panic: %v; stack: %s", p, stack)
			level.Error(log).Log("msg", err)
		}
	}()
	ctxWithTx := context.WithValue(ctx, datasource.TxFieldName, tx)
	accounts, err := p.accountRep.GetForUpdate(ctxWithTx, []string{sender, recipient}, currency)
	if err != nil {
		tx.Rollback()
		level.Error(log).Log("msg", "get for update error "+err.Error())
		return err
	}
	if len(accounts) != 2 {
		tx.Rollback()
		err := fmt.Errorf("expect 2 accounts but found %d: %v", len(accounts), accounts)
		level.Error(log).Log("msg", err)
		return err
	}
	accountFunc := func(owner string) *account.Account {
		for _, a := range accounts {
			if a.Owner == owner {
				return &a
			}
		}
		return nil
	}
	accountFrom := accountFunc(sender)
	accountTo := accountFunc(recipient)
	if accountFrom == nil || accountTo == nil {
		tx.Rollback()
		err := fmt.Errorf("didn't find sender's (%s) or recipient's (%s) account for currency %s: have %v",
			sender, recipient, currency, accounts)
		level.Error(log).Log("msg", err)
		return err
	}
	if accountFrom.Balance < amount {
		tx.Rollback()
		err := fmt.Errorf("unsufficient funds: want to transfer %d but have only %d", amount, accountFrom.Balance)
		level.Error(log).Log("msg", err)
		return err
	}

	accountFrom.Balance -= amount
	updatedRows, err := p.accountRep.UpdateBalance(ctxWithTx, accountFrom.ID, accountFrom.Balance)
	if err != nil {
		tx.Rollback()
		level.Error(log).Log("msg", "update balance error "+err.Error())
		return err
	}
	if updatedRows != 1 {
		tx.Rollback()
		err := fmt.Errorf("count of updated sender's accounts should be 1 but get %d", updatedRows)
		level.Error(log).Log("msg", err)
		return err
	}

	accountTo.Balance += amount
	updatedRows, err = p.accountRep.UpdateBalance(ctxWithTx, accountTo.ID, accountTo.Balance)
	if err != nil {
		tx.Rollback()
		level.Error(log).Log("msg", err)
		return err
	}
	if updatedRows != 1 {
		tx.Rollback()
		err := fmt.Errorf("count of updated recipients's accounts should be 1 but get %d", updatedRows)
		level.Error(log).Log("msg", err)
		return err
	}
	payment := &Payment{AccountFrom: *accountFrom,
		AccountTo: *accountTo,
		Amount:    amount,
		Currency:  accountFrom.Currency,
	}
	err = p.paymentRep.Store(ctxWithTx, payment)
	if err != nil {
		tx.Rollback()
		level.Error(log).Log("msg", "store payment error "+err.Error())
		return err
	}

	return tx.Commit()
}
