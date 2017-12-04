package visitservice

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

//KINDNAME is the table name to store the values
const KINDNAME = "Vist"

//NAMESPACENAME is the Namespace
const NAMESPACENAME = "-kashyak-"

//VisitEntiry is struct to hold the vist details
type VisitEntiry struct {
	ID            int64  `datastore:"-"`
	UserName      string `json:"username,omitempty"`
	Height        string `json:"height,omitempty"`
	Weight        string `json:"weight,omitempty"`
	Temperature   string `json:"temperature,omitempty"`
	BloodPressure string `json:"bloodpressure,omitempty"`
	DoctorNote    string `json:"doctornote,omitempty"`
	PatientNote   string `json:"patientnote,omitempty"`
	NurseNote     string `json:"nursenote,omitempty"`
}

//SuccessResponse store response
type SuccessResponse struct {
	//	visit   VisitEntiry `json:"entity"`
	ID      int64  `json:"Id"`
	Message string `json:"message"`
}

func init() {
	http.HandleFunc("/api/getallvisits/", restHandler)
	http.HandleFunc("/api/postavisit/", restHandler)
	http.HandleFunc("/api/deleteavisit/", restHandler)
}

func restHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		getallvisitshandler(w, r)
		return
	case "POST":
		putavisthandler(w, r)
		return
	case "DELETE":
		deletevisithandler(w, r)
		return
	default:
		//respondErr(w, r, http.StatusNotFound, "is not supported HTTP methods")
	}
}

func getallvisitshandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		panic(err)
	}

	x := strings.Split(r.URL.Path, "/")[3]
	log.Printf("%#v Getting values url - x ", x)
	y, err := strconv.Atoi(x)
	/*if err != nil {
		panic(err)
	}*/
	log.Printf("%#v Getting values url - y ", y)

	if y > 0 {
		// get one
		visit := &VisitEntiry{}
		key := datastore.NewKey(ctx, KINDNAME, "", int64(y), nil)
		if err := datastore.Get(ctx, key, visit); err != nil {
			panic(err)
		}

		visit.ID = key.IntID()
		if err := json.NewEncoder(w).Encode(visit); err != nil {
			panic(err)
		}
		return
	}

	// get all
	visitList := []*VisitEntiry{}
	q := datastore.NewQuery(KINDNAME)
	keys, err := q.GetAll(ctx, &visitList)
	if err != nil {
		panic(err)
	}

	for i, v := range visitList {
		v.ID = keys[i].IntID()
	}
	if err := json.NewEncoder(w).Encode(visitList); err != nil {
		panic(err)
	}
}
func putavisthandler(w http.ResponseWriter, r *http.Request) {

	visit := &VisitEntiry{}
	if err := json.NewDecoder(r.Body).Decode(visit); err != nil {
		panic(err)
	}

	keys := make([]*datastore.Key, 1)
	ctx := appengine.NewContext(r)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		panic(err)
	}

	x := strings.Split(r.URL.Path, "/")[3]
	log.Printf("%#v Getting values url - x ", x)
	y, err := strconv.Atoi(x)
	log.Printf("%#v Getting values url - y ", y)
	if err != nil {
		keys[0] = datastore.NewIncompleteKey(ctx, KINDNAME, nil)
	} else {
		keys[0] = datastore.NewKey(ctx, KINDNAME, "", int64(y), nil)
	}

	k, err := datastore.Put(ctx, keys[0], visit)
	if err != nil {
		panic(err)
	}

	visit.ID = k.IntID()
	json.NewEncoder(w).Encode(visit)
}

func deletevisithandler(w http.ResponseWriter, r *http.Request) {
	var visitslist []VisitEntiry
	var ctx context.Context
	ctx = appengine.NewContext(r)
	keys := make([]*datastore.Key, 1)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		return
	}
	x := strings.Split(r.URL.Path, "/")[3]
	log.Printf("%#v Getting values url - x ", x)
	y, err := strconv.Atoi(x)
	log.Printf("%#v Getting values url - y ", y)
	if y > 0 {

		keys[0] = datastore.NewKey(ctx, KINDNAME, "", int64(y), nil)
		option := &datastore.TransactionOptions{XG: true}
		err = datastore.RunInTransaction(ctx, func(c context.Context) error {
			return datastore.DeleteMulti(c, keys)
		}, option)
	} else {
		q := datastore.NewQuery(KINDNAME)
		keys, _ = q.GetAll(ctx, &visitslist)
		option := &datastore.TransactionOptions{XG: true}
		err = datastore.RunInTransaction(ctx, func(c context.Context) error {
			return datastore.DeleteMulti(c, keys)
		}, option)
	}

}
