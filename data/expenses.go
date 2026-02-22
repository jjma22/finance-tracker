package data

import "errors"

type Expense struct {
	ID    int
	Name  string
	Type  string
	Price float32
	Date  string
}

type Expenses []*Expense

func NewExpense(e *Expense) error {

	if e == nil {
		return errors.New("Enpense cannot be nil")
	}

	if e.Price <= 0 {
		return errors.New("Price has to be greater than 1")
	}



	e.ID = (len(MonthlyExpenses) + 1)
	// Could add some verification on data format
	MonthlyExpenses = append(MonthlyExpenses, e)
	return nil
}

var MonthlyExpenses = Expenses{}

// Module to search if field value already exisits
func Expenses SearchFields(f string, v any) bool, err {

	var currentValues []any
	for _, expense := range MonthlyExpenses {
		currentValues = append(currentNames, expense.f)
	}

	if currentValues.Contains(v) {
		return true
	}
	retur false
}