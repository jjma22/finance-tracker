package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/jjma22/finance-tracker/internal/data"
	"github.com/jjma22/finance-tracker/internal/database"
)

// Define type finance server
type financeServer struct {
	l *slog.Logger
	// add in a validator here so a new one is not created in memory for each request
	v *validator.Validate
}

// Function to return new type finance server with logger
func FinanceNewServer(l *slog.Logger) *financeServer {
	return &financeServer{l, validator.New()}
}

// Function on the financeserver to convert request body to Expense
func (f *financeServer) ExpenseFromJSON(r *http.Request) (error, *data.Expense) {
	var e data.Expense
	// Read request into e
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		f.l.Error(err.Error())
		return err, nil
	}
	fmt.Println(e)
	return nil, &e
}

// Handler to return all expenses
func (f *financeServer) GetExpenses(rw http.ResponseWriter, r *http.Request) {
	f.l.Info("Getting expenses")
	ge, err := database.GetExpenses()
	resp, err := json.Marshal(ge)
	if err != nil {
		f.l.Error("Error getting expenses", "error", err)
		http.Error(rw, "Unable to retive expenses", http.StatusInternalServerError)
	}
	rw.Write(resp)

}

// Handler to return specific expenses
func (f *financeServer) GetExpense(rw http.ResponseWriter, r *http.Request) {
	f.l.Info("Getting expenses")

	// Take id from URL path
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		f.l.Error("Error getting id from path value", "error", err)
		http.Error(rw, "Invalid request", http.StatusBadRequest)
	}

	// Get expense from database via id
	exp, err := database.GetExpense(id)
	if err != nil {
		f.l.Error("Error retrieving expense", "error", err)
		http.Error(rw, "Could not retrieve expense", http.StatusInternalServerError)
	}

	// convert expense into byte slice
	resp, err := json.Marshal(exp)
	if err != nil {
		f.l.Error("Error getting expenses", "error", err)
		http.Error(rw, "Unable to retive expenses", http.StatusInternalServerError)
	}
	rw.Write(resp)

}

// Function to add expense to expenses db
func (f *financeServer) AddExpense(rw http.ResponseWriter, r *http.Request) {

	f.l.Info("Adding new expense")

	// Gets expense from middleware forwarded requests
	e := r.Context().Value(Keyexpense{}).(*data.Expense)

	// Set update date and time
	e.DateAdded = time.Now().Truncate(time.Second)
	e.LastUpdate = time.Now().Truncate(time.Second)

	// Add expense into db
	err := database.AddExpense(e)

	if err != nil {
		f.l.Error("Error adding expense", "error", err)
		http.Error(rw, "Failed to add expense", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(201)
}

func (f *financeServer) UpdateExpense(rw http.ResponseWriter, r *http.Request) {
	f.l.Info("Updating expense")
	// Return type expense from request body
	err, exp := f.ExpenseFromJSON(r)

	if err != nil {
		slog.Error("Issue decoding request body -", "error", err)
		http.Error(rw, "Issue decoding request", http.StatusInternalServerError)
	}

	//Convert id from URL path to int
	exp.ID, _ = strconv.Atoi(r.PathValue("id"))

	// Update expense in db
	err = database.UpdateExpense(exp)
	if err != nil {
		f.l.Error(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

}

// Function to delete expense from db
func (f *financeServer) DeleteExpense(rw http.ResponseWriter, r *http.Request) {

	// Get id from URL path
	ID, _ := strconv.Atoi(r.PathValue("id"))

	// DB query to delete expeense
	// if no rows are changed, rows = 0
	rows, err := database.DeleteExpense(ID)
	if err != nil {
		f.l.Error(err.Error())
		http.Error(rw, "Failed to delete expense", http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		f.l.Error("Cannot delete, ID does not exisit")
		http.Error(rw, "ID not found, cannot delete", http.StatusNotFound)
		return
	}
	f.l.Info("Expense delete", "ID", ID)
}

// Function to return total current expenses
func (f *financeServer) GetTotalExpense(rw http.ResponseWriter, r *http.Request) {
	f.l.Info("Getting total expenses")

	t, err := database.GetTotal()
	if err != nil {
		f.l.Error("Error getting total expenses")
		http.Error(rw, "Errror getting total expenses", http.StatusInternalServerError)
	}

	// en is new encoder that writes to rw Response writer
	en := json.NewEncoder(rw)
	// Writes JSON encoded of t to stream
	en.Encode(t)

}

// Middleware

// Key for context in middleware
type Keyexpense struct{}

// Middleware to validate new expense
func (f *financeServer) MiddleWareValidateExpense(next http.Handler) http.Handler {
	// Annonymous function to validate expense before passing request onto next handler
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		err, expense := f.ExpenseFromJSON(r)
		if err != nil {
			f.l.Error("MW - Error deserialzing product")
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// Function to validate expeense
		err = expense.Validate()
		if err != nil {
			f.l.Error("MW - Error validating expense")
			http.Error(
				rw,
				fmt.Sprintf("Error validating expense: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// Add key with expense to context
		// May be incorrect, original request should already contain request so should it be decoded again??
		ctx := context.WithValue(r.Context(), Keyexpense{}, expense)
		r = r.WithContext(ctx)

		// calls next handler passed in as next, currently AddExpense
		next.ServeHTTP(rw, r)
	})
}
