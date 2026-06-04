package data

import "testing"

func TestChecksValidatation(t *testing.T) {
	e := &Expense{
		Name:  "Bills",
		Price: 60,
		SKU:   "abs-sbc-sgv",
	}

	err := e.Validate()

	if err != nil {
		t.Fatal(err)
	}
}

func TestNameValidation(t *testing.T) {
	t.Run("Test if expense validation blocks expense without name", func(t *testing.T) {
		e := &Expense{
			Price: 60,
			SKU:   "abs-sbc-sgv",
		}

		err := e.Validate()
		//Fail if no error is returned
		if err == nil {
			t.Fatal(err)
		}

	})
}

func TestPriceValidationNotSet(t *testing.T) {
	t.Run("Test if expense validation blocks expense without price", func(t *testing.T) {
		e := &Expense{
			Name: "Test",
			SKU:  "abs-sbc-sgv",
		}

		err := e.Validate()
		//Fail if no error is returned
		if err == nil {
			t.Fatal(err)
		}

	})
}

func TestPriceValidationSetZero(t *testing.T) {
	t.Run("Test if expense validation blocks expense with price set to 0", func(t *testing.T) {
		e := &Expense{
			Name:  "Test",
			Price: 0,
			SKU:   "abs-sbc-sgv",
		}

		err := e.Validate()
		//Fail if no error is returned
		if err == nil {
			t.Fatal(err)
		}

	})
}

func TestPriceValidationSetNegative(t *testing.T) {
	t.Run("Test if expense validation blocks expense with price set to negative", func(t *testing.T) {
		e := &Expense{
			Name:  "Test",
			Price: -1000,
			SKU:   "abs-sbc-sgv",
		}

		err := e.Validate()
		//Fail if no error is returned
		if err == nil {
			t.Fatal(err)
		}

	})
}

func TestPriceValidationSKU(t *testing.T) {
	t.Run("Test if expense validation blocks expense without SKU", func(t *testing.T) {
		e := &Expense{
			Name:  "Test",
			Price: 1000,
		}

		err := e.Validate()
		//Fail if no error is returned
		if err == nil {
			t.Fatal(err)
		}

	})
}

// Accepted SKU format is `[a-z]+-[a-z]+-[a-z]+`
func TestPriceValidationInvalidSKU(t *testing.T) {
	t.Run("Test if expense validation blocks expense with invalid SKU", func(t *testing.T) {
		e := &Expense{
			Name:  "Test",
			Price: 1000,
			SKU:   "aac-fge",
		}

		err := e.Validate()
		//Fail if no error is returned
		if err == nil {
			t.Fatal(err)
		}

	})
}

func TestPriceValidationValidExpense(t *testing.T) {
	t.Run("Test if expense validation blocks expense with invalid SKU", func(t *testing.T) {
		e := &Expense{
			Name:  "Test",
			Price: 1000,
			SKU:   "aac-fge-hjj",
		}

		err := e.Validate()
		//Fail if no error is returned
		if err != nil {
			t.Fatal(err)
		}

	})
}
