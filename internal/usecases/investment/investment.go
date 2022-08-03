package investment

import (
	"context"

	"github.com/powerman/structlog"
)

type Usecase interface {
	ListInvestments(ctx context.Context, req ListInvestmentRequest) (ListInvestmentResponse, error)
	CreateInvestment(ctx context.Context, ticker string, price, quantity int64) error
	UpdateInvestment(ctx context.Context, id, price, quantity int64) error
	DeleteInvestment(ctx context.Context, id int64) error
}

type Portfolio struct {
	repo   Repo
	logger *structlog.Logger
}

func New(repo Repo, logger *structlog.Logger) *Portfolio {
	return &Portfolio{repo, logger}
}

func (p *Portfolio) ListInvestments(ctx context.Context, req ListInvestmentRequest) (ListInvestmentResponse, error) {
	return p.repo.ListInvestment(ctx, req)
}

func (p *Portfolio) CreateInvestment(ctx context.Context, ticker string, price, quantity int64) error {
	return p.repo.CreateInvestment(ctx, ticker, price, quantity)
}

func (p *Portfolio) UpdateInvestment(ctx context.Context, id, price, quantity int64) error {
	return p.repo.UpdateInvestment(ctx, id, price, quantity)
}

func (p *Portfolio) DeleteInvestment(ctx context.Context, id int64) error {
	return p.repo.DeleteInvestment(ctx, id)
}
