package hellodatastore

import (
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	http.Handle("/", &templateHandler{filename: "index.html"})
	//tpl = template.Must(template.ParseGlob("templates/html/*.html"))
	//tpl = template.Must(template.ParseGlob("templates/*.html"))
	tpl = template.Must(template.ParseGlob("templates/*"))
	http.HandleFunc("/api/", restHandler)
	//http.HandleFunc("/home", home)
	log.Printf("%#v", "Hi init() ")

}

/*
func home(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)

	if len(session.Value) > 0 {
		var sd sessionData
		json.Unmarshal(session.Value, &sd)
		sd.LoggedIn = true
		tpl.ExecuteTemplate(w, "home.html", nil)
	}
}*/
func restHandler(w http.ResponseWriter, r *http.Request) {

	p := NewParam(r.URL)
	/*if p.Kind != "user" {
		respondErr(w, r, http.StatusBadRequest, "does not correspond to the type of non-Expense")
	}*/
	ed := NewUserDatasotre(r)

	if err := ed.build(p.UName, p.Fname, p.LName, p.Email, p.Pass); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err.Error())
		log.Printf("%#v build ", http.StatusInternalServerError)
		log.Printf("%#v build ", err.Error())
	}
	//serveTemplate(w, r, "index.html")
	switch r.Method {
	case "GET":
		handleGet(ed, w, r)
		return
	case "PUT":
		handlePut(ed, w, r)
		return
	case "DELETE":
		handleDelete(ed, w, r)
		return
	case "POST":
		handlePost(ed, w, r)
		return
	default:
		respondErr(w, r, http.StatusNotFound, "is not supported HTTP methods")
	}
}

//SuccessResponse is used to send sucess message
type SuccessResponse struct {
	User    UserEntiry `json:"entity"`
	Message string     `json:"message"`
}

func handleGet(ed *UserDatasotre, w http.ResponseWriter, r *http.Request) {
	/*if err := ed.Get(w, r); err != nil {
		respondErr(w, r, http.StatusBadRequest, err.Error())
		return
	}
	message := "「" + ed.User.FirstName + "amount of money is " + ed.User.UserName + "It is yen"
	respond(w, r, http.StatusOK, SuccessResponse{*ed.User, message})
	log.Printf("%#v message ", message)
	log.Printf("%#v message ", http.StatusOK)*/
}

func handlePut(ed *UserDatasotre, w http.ResponseWriter, r *http.Request) {
	if err := ed.createUser(ed, w, r); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	message := "「" + ed.User.FirstName + "made a registration of"
	respond(w, r, http.StatusOK, SuccessResponse{*ed.User, message})
}

func handleDelete(ed *UserDatasotre, w http.ResponseWriter, r *http.Request) {
	/*	if err := ed.Delete(); err != nil {
			respondErr(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		message := "「" + ed.User.FirstName + "gave a deletion of"
		respond(w, r, http.StatusOK, SuccessResponse{*ed.User, message})*/
}

func handlePost(ed *UserDatasotre, w http.ResponseWriter, r *http.Request) {
	ed.login(ed, w, r)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
	/*
		if err := ed.login(ed, w, r); err != nil {
			respondErr(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		message := "「" + ed.User.FirstName + "Logged INNNN"
		respond(w, r, http.StatusOK, SuccessResponse{*ed.User, message})
	*/
}
