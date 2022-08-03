package investment

import "context"

type Repo interface {
	ListInvestment(ctx context.Context, req ListInvestmentRequest) (ListInvestmentResponse, error)
	CreateInvestment(ctx context.Context, ticker string, price, quantity int64) error
	UpdateInvestment(ctx context.Context, id, price, quantity int64) error
	DeleteInvestment(ctx context.Context, id int64) error
}
