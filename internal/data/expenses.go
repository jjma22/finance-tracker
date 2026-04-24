package data

import (
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type Expense struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required"`
	// Type  string `json:"type"`
	Price      float32   `json:"price" validate:"gt=0"`
	SKU        string    `json:"sku" validate:"required,sku"`
	DateAdded  time.Time `json:"-"`
	LastUpdate time.Time `json:"-"`
}

type Expenses []*Expense

func (e *Expense) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(e)

}

func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
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
