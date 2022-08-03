package investment

import "time"

type ListInvestmentRequest struct {
}

type ListInvestmentResponse struct {
	Investments []Investment
}

type Investment struct {
	ID        int64     `db:"id" json:"id"`
	Ticker    string    `db:"ticker" json:"ticker"`
	Quantity  int64     `db:"quantity" json:"quantity"`
	Price     int64     `db:"price" json:"price"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
