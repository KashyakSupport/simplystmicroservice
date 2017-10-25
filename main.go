package simplystmicroservice

import "net/http"

func init() {
	http.Handle("/", &templateHandler{filename: "index.html"})
	http.HandleFunc("/api/", restHandler)
}

func restHandler(w http.ResponseWriter, r *http.Request) {

	p := NewParam(r.URL)

	if p.Kind != "user" {
		respondErr(w, r, http.StatusBadRequest, "types other than user not supportted")

	}
	if r.Method != "PUT" && !p.HasValue() {
		p.Value = "0"
	}

	/*price, err := strconv.Atoi(p.Value)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, err.Error())
		return
	}*/

	ud := NewUserDatastore(r)

	switch r.Method {
	case "GET":
		handleGet(ud, w, r)
		return
	case "PUT":
		handlePut(ud, w, r)
		return
	default:
		respondErr(w, r, http.StatusNotFound, "Error from defalut case")

	}

}

// SuccessResponse is the best response
type SuccessResponse struct {
	User    UserEntiry `json:"entity"`
	Message string     `json:"message"`
}

func handleGet(ud *UserDatastore, w http.ResponseWriter, r *http.Request) {
	if err := ud.Get(); err != nil {
		respondErr(w, r, http.StatusBadRequest, err.Error())
		return
	}
}

func handlePut(ud *UserDatastore, w http.ResponseWriter, r *http.Request) {

	if err := ud.Put(); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	message := "Hi" + ud.User.FirstName + "Successfully Created."
	respond(w, r, http.StatusOK, SuccessResponse{*ud.User, message})
}
