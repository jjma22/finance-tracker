package data

type Budget struct {
	Budget int `json:"budget"`
}

var MonthlyBudget = &Budget{
	Budget: 0,
}

func UpdateBudget(i int) error {

	MonthlyBudget.Budget = i
	return nil
}