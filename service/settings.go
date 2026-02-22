package service

import "main.go/data"

func GetBudget() (*data.Budget, error) {
	return data.MonthlyBudget, nil
}

