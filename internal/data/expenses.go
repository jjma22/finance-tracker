package data

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

type Expense struct {
	ID    int
	Name  string
	Type  string
	Price float32
	DateAdded  string
	LastUpdate string
}

type Expenses []*Expense

//Simply func but easier to rewrite
//once database is added
func GetExpenses() Expenses{
	return MonthlyExpenses
}

func NewExpense(e *Expense) error {

	if e == nil {
		return errors.New("Enpense cannot be nil")
	}

	if e.Price <= 0 {
		return errors.New("Price has to be greater than 1")
	}



	e.ID = (len(MonthlyExpenses) + 1)

	b, err := e.SearchFields("Name", e.Name)

	if err != nil {
		return errors.New("Checking invalid field")
	}
	if b == true {
		return errors.New("Item Name already exists in expenses")
	}

	//Set DateAdded
	e.DateAdded = time.Now().Truncate(time.Second).Format("2006-01-02 15:04:05")
	// Could add some verification on data format
	MonthlyExpenses = append(MonthlyExpenses, e)
	return nil
}

var MonthlyExpenses = Expenses{}

// Module to search if field value already exisits
// This will need rewriting when applicaiton uses database
func (Expense) SearchFields(f any, v any) (bool, error) {
	// Create array of values from specific field
	var currentValues []any
	for _, expense := range MonthlyExpenses {
		switch f{
		case "Name":
			currentValues = append(currentValues, expense.Name)
		case "ID":
			currentValues = append(currentValues, expense.ID)
		default:
			return false, errors.New("Invalid field")

		}
	}

	// return true if slice contains values
	if slices.Contains(currentValues, v) {
		return true, nil
	} else {
		return false, nil
	}
}

func UpdateExpense(e *Expense) error {
	v, err := e.SearchFields("ID", e.ID)
	if err != nil {
		return errors.New("Error searching for ID")
	}
	if v == false {
		return errors.New("ID does not exist")
	}
	fmt.Println("Expense found, updating")
	//MonthlyExpenses[(e.ID - 1)].Price = e.Price
	MonthlyExpenses[(e.ID - 1)].Price = e.Price
	MonthlyExpenses[(e.ID - 1)].LastUpdate = time.Now().Truncate(time.Second).Format("2006-01-02 15:04:05")
	return nil

}

func DeleteExpense(i int) error {
	var e Expense
	e.ID = i
	v, err := e.SearchFields("ID", i)
	if err != nil {
		return errors.New("Error searching for ID")
	}
	if v == false {
		return errors.New("ID does not exist")
	}

	MonthlyExpenses = slices.Delete(MonthlyExpenses, (i-1), i)
	return nil

}

func GetTotal() (float32, error) {
	var t float32
	for _, e := range MonthlyExpenses {
		t += e.Price
	}
	return t, nil
}