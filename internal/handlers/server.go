package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jjma22/finance-tracker.git/internal/data"
	"github.com/jjma22/finance-tracker.git/internal/database"
	"github.com/jjma22/finance-tracker.git/internal/service"
)

type financeServer struct {
	l *slog.Logger
}

func FinanceNewServer(l *slog.Logger) *financeServer{
	return &financeServer{l}
}


func (f*financeServer) GetBudget(rw http.ResponseWriter, r*http.Request) {
	mb, err := service.GetBudget()

	if err != nil {
		f.l.Error("Error calling data.GetBudget")
		http.Error(rw, "Unable to get monthly budget", http.StatusInternalServerError)
	}

	d, err := json.Marshal(mb)
	if err != nil {
		f.l.Error("Unable to marhshal budget")
		http.Error(rw, "Unable to get monthly budget", http.StatusInternalServerError)
	}

	f.l.Info("Returning current budget")
	rw.Write(d)

}

func (f*financeServer) UpdateBudget(rw http.ResponseWriter, r*http.Request) {
	var mb data.Budget
	f.l.Info("Updating budget")
	err := json.NewDecoder(r.Body).Decode(&mb)
	if err != nil {
		f.l.Error(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	err = data.UpdateBudget(mb.Budget)

	if err != nil {
		f.l.Error(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	f.l.Info("Budget updated")
	

}

func (f*financeServer) GetExpenses(rw http.ResponseWriter, r*http.Request) {
	f.l.Info("Getting expenses")
	ge := data.GetExpenses()
	// ge := database.GetExpense
	resp, err := json.Marshal(ge)
	if err != nil {
		f.l.Error("Error getting expenses", "error", err)
		http.Error(rw, "Unable to retive expenses", http.StatusInternalServerError)
	}
	rw.Write(resp)


}
 
func (f*financeServer) ExpenseFromJSON(r*http.Request) (error, *data.Expense) {
	var e data.Expense
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		f.l.Error(err.Error())
		return err, nil
	}
	fmt.Println(e)
	return nil, &e
}

func (f*financeServer) AddExpense(rw http.ResponseWriter, r*http.Request) {
	f.l.Info("Adding new expense")
	e := r.Context().Value(Keyexpense{}).(*data.Expense)

	fmt.Println(e)

	err := data.NewExpense(e)

	if err != nil {
		f.l.Error("Error adding expense", "error", err)
		http.Error(rw, "Failed to add expense", http.StatusInternalServerError)
	}


}

func (f*financeServer) UpdateExpense( rw http.ResponseWriter, r*http.Request) {
	f.l.Info("Updating expense")
	
	var exp data.Expense
	err := json.NewDecoder(r.Body).Decode(&exp)
	if err != nil {
		f.l.Error("Error unmarshalling request")
		http.Error(rw, "Unable to unmarshall request", http.StatusInternalServerError)
		return
	}

	exp.ID, _ = strconv.Atoi(r.PathValue("id"))
	err = data.UpdateExpense(&exp)
	if err != nil {
		f.l.Error(err.Error())
	}
	
}

func (f*financeServer)DeleteExpense(rw http.ResponseWriter, r * http.Request) {

		ID, _ := strconv.Atoi(r.PathValue("id"))
		err := data.DeleteExpense(ID)
		if err != nil {
			f.l.Error(err.Error())
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		f.l.Info("Expense delete", "ID", ID)
}

func (f*financeServer)GetTotalExpense(rw http.ResponseWriter, r *http.Request) {
	f.l.Info("Getting total expenses")

	t, err := database.GetTotal()
	if err != nil {
		f.l.Error("Error getting total expenses")
		http.Error(rw, "Errror getting total expenses", http.StatusInternalServerError)
	}
	//rw.Write(Float32ToByte(float32(val)))
	en := json.NewEncoder(rw)
	en.Encode(t)


}

type Keyexpense struct{}

func (f*financeServer) MiddleWareValidateExpense(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r*http.Request) {
		err, expense := f.ExpenseFromJSON(r)
		if err != nil {
			f.l.Error("MW - Error deserialzing product")
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

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
		fmt.Println(expense)
		ctx := context.WithValue(r.Context(), Keyexpense{}, expense)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}

