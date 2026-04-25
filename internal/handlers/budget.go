package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jjma22/finance-tracker.git/internal/data"
	"github.com/jjma22/finance-tracker.git/internal/database"
)

func (f *financeServer) BudgetFromJSON(r *http.Request) (*data.Budget, error) {
	var b data.Budget
	// Read request into b
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		f.l.Error(err.Error())
		return nil, err
	}
	return &b, nil
}

func (f *financeServer) SetBudget(rw http.ResponseWriter, r *http.Request) {
	b := r.Context().Value(Budget{}).(*data.Budget)
	err := database.SetBudget(b.Budget)

	if err != nil {
		f.l.Error("Error setting new budget in databse", "err", err)
		http.Error(rw, "Error occured when adding budget", http.StatusInternalServerError)
	}

}

// Handler to return monthly budget
func (f *financeServer) GetBudget(rw http.ResponseWriter, r *http.Request) {

	// Take id from URL path
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		f.l.Error("Error getting id from path value", "error", err)
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	b, err := database.GetBudget(id)

	if err != nil {
		f.l.Error("Error returning budget", "error", err)
		http.Error(rw, "Unable to get monthly budget", http.StatusInternalServerError)
		return
	}

	d, err := json.Marshal(b)
	if err != nil {
		f.l.Error("Unable to marhshal budget")
		http.Error(rw, "Unable to get monthly budget", http.StatusInternalServerError)
		return
	}

	f.l.Info("Returning current budget")
	rw.Write(d)

}

// Handler to update monthly budget
func (f *financeServer) UpdateBudget(rw http.ResponseWriter, r *http.Request) {

	// Take id from URL path
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		f.l.Error("Error getting id from path value", "error", err)
		http.Error(rw, "Invalid request", http.StatusBadRequest)
	}

	// Decode request body into mb
	b, err := f.BudgetFromJSON(r)
	if err != nil {
		f.l.Error("Error decoding request at UpdateBudget", "error", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update monthly budget
	err = database.UpdateBudget(id, b.Budget)

	if err != nil {
		f.l.Error(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	f.l.Info("Budget updated")

}

type Budget struct{}

// Middleware to validate new expense
func (f *financeServer) MiddleWareValidateBudget(next http.Handler) http.Handler {
	// Annonymous function to validate budget before passing request onto next handler
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		b, err := f.BudgetFromJSON(r)
		if err != nil {
			f.l.Error("MW - Error deserialzing product")
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		f.l.Info("Validating new budget")
		// Function to validate expeense
		err = f.v.Struct(b)
		if err != nil {
			f.l.Error("MW - Error validating budget", "error", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating budget: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// Add key with expense to context
		// May be incorrect, original request should already contain request so should it be decoded again??
		f.l.Info("Passing budget to next hanler")
		ctx := context.WithValue(r.Context(), Budget{}, b)
		r = r.WithContext(ctx)

		// calls next handler passed in as next, currently AddExpense
		next.ServeHTTP(rw, r)
	})
}
