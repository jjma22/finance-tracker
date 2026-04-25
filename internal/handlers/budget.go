package handlers

import (
	"encoding/json"
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
	f.l.Info("Setting new budget")
	b, err := f.BudgetFromJSON(r)
	if err != nil {
		f.l.Error("Error decoding budget", "error", err)
		http.Error(rw, "Error occured setting budget", http.StatusInternalServerError)
	}

	err = database.SetBudget(b.Budget)

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
