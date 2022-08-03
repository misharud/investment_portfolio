package postgres

import (
	"context"

	"github.com/misharud/investment_portfolio/internal/usecases/investment"

	"github.com/pkg/errors"
)

func (r *Repo) ListInvestment(ctx context.Context, req investment.ListInvestmentRequest) (investment.ListInvestmentResponse, error) {
	q := `select id, ticker, quantity, price, created_at from investment_history`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return investment.ListInvestmentResponse{}, errors.WithStack(err)
	}
	defer rows.Close()

	var resp investment.ListInvestmentResponse
	for rows.Next() {
		var p investment.Investment
		err = rows.Scan(&p.ID, &p.Ticker, &p.Quantity, &p.Price, &p.CreatedAt)
		if err != nil {
			return investment.ListInvestmentResponse{}, errors.WithStack(err)
		}
		resp.Investments = append(resp.Investments, p)
	}
	err = rows.Err()
	if err != nil {
		return investment.ListInvestmentResponse{}, errors.WithStack(err)
	}
	return resp, nil
}

func (r *Repo) CreateInvestment(ctx context.Context, ticker string, price, quantity int64) error {
	_, err := r.db.ExecContext(ctx,
		`
		insert into investment_history(ticker, price, quantity)
		values(:ticker, :price, :quantity)
	`, ticker, price, quantity)
	return err
}

func (r *Repo) UpdateInvestment(ctx context.Context, id, price, quantity int64) error {
	_, err := r.db.ExecContext(ctx, `
		update investment_history set price = $1, quantity = $2 where id = $3
	`, price, quantity, id)
	return err
}

func (r *Repo) DeleteInvestment(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "delete from investment_portfolio where id = $1", id)
	return err
}
