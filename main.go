package hellodatastore

import (
	"log"
	"net/http"
)

func init() {
	http.Handle("/", &templateHandler{filename: "index.html"})
	http.HandleFunc("/api/", restHandler)
}

func restHandler(w http.ResponseWriter, r *http.Request) {

	p := NewParam(r.URL)
	log.Printf("%#v", r.URL)
	log.Printf("%#v", p)
	if p.Kind != "expense" {
		log.Printf("%#v", p.Kind)
		respondErr(w, r, http.StatusBadRequest, "does not correspond to the type of non-Expense")
	}

	/*if r.Method != "PUT" && !p.HasValue() {
		p.Value = "0"
	}*/

	/*price, err := strconv.Atoi(p.Value)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, err.Error())
		return
	}*/
	log.Printf("%#v", r)
	ed := NewExpenseDatasotre(r)
	log.Printf("%#v", ed)
	if err := ed.build(p.Key, p.LName, p.UName, p.Pass, p.Email); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err.Error())
	}

	switch r.Method {
	/*	case "GET":
		handleGet(ed, w, r)
		return*/
	case "PUT":
		handlePut(ed, w, r)
		return
	/*case "DELETE":
	handleDelete(ed, w, r)
	return*/
	default:
		respondErr(w, r, http.StatusNotFound, "is not supported HTTP methods")
	}
}

type SuccessResponse struct {
	Expense ExpenseEntiry `json:"entity"`
	Message string        `json:"message"`
}

/*func handleGet(ed *ExpenseDatasotre, w http.ResponseWriter, r *http.Request) {
	if err := ed.Get(); err != nil {
		respondErr(w, r, http.StatusBadRequest, err.Error())
		return
	}
	message := "「" + ed.Expense.Name + "amount of money is " + strconv.Itoa(ed.Expense.Price) + "It is yen"
	respond(w, r, http.StatusOK, SuccessResponse{*ed.Expense, message})
}*/

func handlePut(ed *ExpenseDatasotre, w http.ResponseWriter, r *http.Request) {
	if err := ed.Put(); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	message := "「" + ed.Expense.Name + "made a registration of"
	respond(w, r, http.StatusOK, SuccessResponse{*ed.Expense, message})
}

/*func handleDelete(ed *ExpenseDatasotre, w http.ResponseWriter, r *http.Request) {
	if err := ed.Delete(); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	message := "「" + ed.Expense.Name + "gave a deletion of"
	respond(w, r, http.StatusOK, SuccessResponse{*ed.Expense, message})
}*/
