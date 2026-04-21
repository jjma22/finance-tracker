package data

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"time"

	"github.com/go-playground/validator"
)

type Expense struct {
	ID    int `json:"id"`
	Name  string `json:"name" validate:"required"`
	// Type  string `json:"type"`
	Price float32 `json:"price" validate:"gt=0"`
	SKU string `json:"sku" validate:"required,sku"`
	DateAdded  time.Time `json:"-"`
	LastUpdate time.Time `json:"-"`
}

type Expenses []*Expense

//Simply func but easier to rewrite
//once database is added
func GetExpenses() Expenses{
	return MonthlyExpenses
}

func (e*Expense) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(e)

}

func validateSKU(fl validator.FieldLevel) bool {
	re :=  regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	if fl.Field().String() == "invalid" {
		return false
	}

	return true
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
	fmt.Println(e.Price)
	MonthlyExpenses[(e.ID - 1)].Price = e.Price
	//MonthlyExpenses[(e.ID - 1)].LastUpdate = time.Now().Truncate(time.Second).Format("2006-01-02 15:04:05")
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

// No longer needed now db implemented
// func GetTotal() (float32, error) {
// 	var t float32
// 	for _, e := range MonthlyExpenses {
// 		t += e.Price
// 	}
// 	return t, nil
// }