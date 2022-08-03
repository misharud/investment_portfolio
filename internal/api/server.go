package api

import (
	"context"
	"log"
	"net/http"
	"runtime"

	"github.com/misharud/investment_portfolio/internal/api/handlers"
	"github.com/misharud/investment_portfolio/internal/provider/postgres"
	"github.com/misharud/investment_portfolio/internal/usecases/investment"
	"github.com/misharud/investment_portfolio/pkg/config"

	"github.com/gorilla/mux"
	"github.com/powerman/structlog"
)

func StartServer() error {
	ctx := context.Background()
	logger := structlog.New(structlog.KeyUnit, "main")

	conf, err := config.Load()
	if err != nil {
		logger.Fatal("config parsing", "err", err)
	}

	runtime.GOMAXPROCS(conf.Server.MaxCPU)
	logger.Println("CPU", "GOMAXPROCS", runtime.GOMAXPROCS(-1))

	db, err := postgres.Open(ctx, &conf.Database, logger)
	if err != nil {
		logger.Fatal("db connection", "err", err)
	}

	err = postgres.MigrationDB(db, conf.Database.MigrationPath)
	if err != nil {
		logger.Fatal("migrations was failure", "err", err)
	}

	repo := postgres.NewRepo(db, logger)

	h := handlers.NewHandler(investment.New(repo, logger), logger)

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/investments", h.ListInvestments).Methods(http.MethodGet)
	r.HandleFunc("/investment", h.CreateInvestment).Methods(http.MethodPost)
	r.HandleFunc("/investment/:id", h.UpdateInvestment).Methods(http.MethodPut)
	r.HandleFunc("/investment/:id", h.DeleteInvestment).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(conf.Server.REST.Addr(), r))

	return nil
}
