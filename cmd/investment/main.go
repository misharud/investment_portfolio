package main

import (
	"os"

	"github.com/misharud/investment_portfolio/internal/api"
)

func main() {
	if err := api.StartServer(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
