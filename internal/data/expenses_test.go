package data

import "testing"

func TestChecksValidatation(t *testing.T) {
	e := &Expense{
		Name: "Bills",
		Price: 60,
		SKU: "abs-sbc-sgv",
	}

	err := e.Validate()

	if err != nil {
		t.Fatal(err)
	}
}

