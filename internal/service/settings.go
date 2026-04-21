package service

import "github.com/jjma22/finance-tracker.git/internal/data"

func GetBudget() (*data.Budget, error) {
	return data.MonthlyBudget, nil
}

