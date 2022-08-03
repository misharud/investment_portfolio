package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/misharud/investment_portfolio/internal/usecases/investment"

	"github.com/gorilla/mux"
	"github.com/powerman/structlog"
)

type Handler struct {
	investmentSvc investment.Usecase
	logger        *structlog.Logger
}

func NewHandler(investmentSvc investment.Usecase, logger *structlog.Logger) *Handler {
	return &Handler{
		investmentSvc: investmentSvc,
		logger:        logger,
	}
}

func (h *Handler) ListInvestments(w http.ResponseWriter, r *http.Request) {
	i, err := h.investmentSvc.ListInvestments(r.Context(), investment.ListInvestmentRequest{})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(i.Investments)
}

func (h *Handler) CreateInvestment(w http.ResponseWriter, r *http.Request) {
	var m investment.Investment

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.PrintErr("CreateInvestment decode", "err", err)
		return
	}

	err := h.investmentSvc.CreateInvestment(r.Context(), m.Ticker, m.Price, m.Quantity)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.PrintErr("CreateInvestment db", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateInvestment(w http.ResponseWriter, r *http.Request) {
	var m investment.Investment
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.PrintErr("UpdateInvestment decode", "err", err)
		return
	}

	err := h.investmentSvc.UpdateInvestment(r.Context(), int64(id), m.Price, m.Quantity)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.PrintErr("UpdateInvestment db", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteInvestment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	err := h.investmentSvc.DeleteInvestment(r.Context(), int64(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.PrintErr("DeleteInvestment db", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
