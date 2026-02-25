package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"main.go/data"
	"main.go/service"
)

type financeServer struct {
	l *log.Logger
}

func FinanceNewServer(l *log.Logger) *financeServer{
	return &financeServer{l}
}

func (f*financeServer) ServeHTTP(rw http.ResponseWriter, r*http.Request) {
	// refactor so cpature is case insensitive
	if r.Method == http.MethodGet && r.URL.Path == "/monthlybudget"{
		f.GetBudget(rw, r)
		return

	}

	if r.Method == http.MethodPut && r.URL.Path == "/monthlybudget"{
		f.UpdateBudget(rw, r)
		return
	}

	if r.Method == http.MethodGet && r.URL.Path == "/expense"{
		f.GetExpenses(rw, r)
		return
	}

	if r.Method == http.MethodPost && r.URL.Path == "/expense"{
		f.AddExpense(rw, r)
		return
	}

	if r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/expense/update/") {
		f.UpdateExpense(rw, r)
	}

	if r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/expense/delete/") {
		f.DeleteExpense(rw, r)
	}
}


func (f*financeServer) GetBudget(rw http.ResponseWriter, r*http.Request) {
	mb, err := service.GetBudget()

	if err != nil {
		f.l.Println("Error calling data.GetBudget")
		http.Error(rw, "Unable to get monthly budget", http.StatusInternalServerError)
	}

	d, err := json.Marshal(mb)
	if err != nil {
		f.l.Println("Unable to marhshal budget")
		http.Error(rw, "Unable to get monthly budget", http.StatusInternalServerError)
	}

	f.l.Println("Returning current budget")
	rw.Write(d)

}

func (f*financeServer) UpdateBudget(rw http.ResponseWriter, r*http.Request) {
	var mb data.Budget
	f.l.Println("Updating budget")
	err := json.NewDecoder(r.Body).Decode(&mb)
	if err != nil {
		f.l.Println(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	err = data.UpdateBudget(mb.Budget)

	if err != nil {
		f.l.Println(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	f.l.Println("Budget updated")
	

}

func (f*financeServer) GetExpenses(rw http.ResponseWriter, r*http.Request) {
	f.l.Println("Getting expenses")
	ge := data.GetExpenses()
	resp, err := json.Marshal(ge)
	if err != nil {
		f.l.Printf("Error getting expenses %s", err)
		http.Error(rw, "Unable to retive expenses", http.StatusInternalServerError)
	}
	rw.Write(resp)


}

func (f*financeServer) AddExpense(rw http.ResponseWriter, r*http.Request) {
	// Read body with expense into ne
	var ne data.Expense
	f.l.Println("Adding new expense")
	err := json.NewDecoder(r.Body).Decode(&ne)
	if err != nil {
		f.l.Println(err)
		http.Error(rw, "Unable to unmarshall json", http.StatusInternalServerError)
	}
	f.l.Println(ne)

	//Add new expense into expenses
	err = data.NewExpense(&ne)

	if err != nil {
		f.l.Printf("Error adding expense: %s", err)
		http.Error(rw, "Failed to add expense", http.StatusInternalServerError)
	}


 }

 func (f*financeServer) UpdateExpense( rw http.ResponseWriter, r*http.Request) {
	f.l.Println("Updating expense")
	
	// Get ID from path
	e := strings.Split(r.URL.Path,`/`)
	f.l.Println(len(e))
	if (len(e) != 4) {
		f.l.Println("Bad request, invalid path")
		http.Error(rw, "Invalid path", http.StatusBadRequest)
	}
	id, err := strconv.Atoi(e[3])
	if err != nil {
		f.l.Println("Invalid request sent to UpdateExpense")
		http.Error(rw, "Bad request", http.StatusBadRequest)
	}

	var exp data.Expense
	err = json.NewDecoder(r.Body).Decode(&exp)
	if err != nil {
		f.l.Println("Error unmarshalling request")
		http.Error(rw, "Unable to unmarshall request", http.StatusInternalServerError)
	}

	exp.ID = id
	err = data.UpdateExpense(&exp)
	if err != nil {
		f.l.Println(err)
	}
	
 }

 func (f*financeServer)DeleteExpense(rw http.ResponseWriter, r * http.Request) {
		// Get ID from path
		e := strings.Split(r.URL.Path,`/`)
		f.l.Println(len(e))
		if (len(e) != 4) {
			f.l.Println("Bad request, invalid path")
			http.Error(rw, "Invalid path", http.StatusBadRequest)
		}
		id, err := strconv.Atoi(e[3])
		if err != nil {
			f.l.Println("Invalid request sent to UpdateExpense")
			http.Error(rw, "Bad request", http.StatusBadRequest)
		}

		err = data.DeleteExpense(id)
 }