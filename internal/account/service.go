package account

import (
	"context"

	"github.com/go-kit/kit/log/level"
	"github.com/sgdr/wallet-service/internal/logger"
)

type Service interface {
	GetAll(ctx context.Context) ([]Account, error)
}

type accountService struct {
	r Repository
}

func NewService(r Repository) Service {
	return &accountService{r: r}
}

func (a accountService) GetAll(ctx context.Context) ([]Account, error) {
	log := logger.FromContext(ctx)
	level.Info(log).Log("msg", "get all accounts")
	return a.r.FindAll(ctx)
}
