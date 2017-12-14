package visitservice

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/codegangsta/negroni"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

//KINDNAME is the table name to store the values
const KINDNAME = "Visit"

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

//User is a strct
type User struct {
	Username string `datastore:"-"`
	Password string `json:"password"`
}

// Response is strcut
type Response struct {
	Data string `json:"data"`
}

// Token is also strcut
type Token struct {
	Token string `json:"token"`
}

const (
	privKeyPath = "keys/app.rsa"
	pubKeyPath  = "keys/app.rsa.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)
var verifyBytes, signBytes []byte

func init() {
	initKeys()
	startServer()
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
func initKeys() {
	var err error

	signBytes, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatalf("Error reading private key: %v", err)
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatalf("Error parsing private key: %v", err)
	}
	verifyBytes, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatalf("Error reading public key: %v", err)
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatalf("Error parsing public key: %v", err)
	}
}

// ProtectedHandler is handler
/*func ProtectedHandler(w http.ResponseWriter, r *http.Request) {

	response := Response{"Gained access to protected resource"}
	JSONResponse(response, w)

}*/

//RegisterHandler is using to insert user into database
func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		panic(err)
	}

	keys := make([]*datastore.Key, 1)
	ctx := appengine.NewContext(r)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		panic(err)
	}

	keys[0] = datastore.NewKey(ctx, "User", user.Username, 0, nil)

	_, err = datastore.Put(ctx, keys[0], user)
	if err != nil {
		panic(err)
	}

	//user.ID = k.IntID()
	json.NewEncoder(w).Encode(user)
}

//LoginHandler is using to check wherther user existed or not
//If User Existed In the datastore send tocken back
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		panic(err)
	}
	ctx := appengine.NewContext(r)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		panic(err)
	}

	// get one
	//user := &User{}
	key := datastore.NewKey(ctx, "User", user.Username, 0, nil)
	err = datastore.Get(ctx, key, user)
	if err != nil {
		// there is an err, there is a NO user
		//fmt.Fprint(w, "false")
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Error logging in")
		fmt.Fprint(w, "Invalid credentials")
		return
	}

	//create a rsa 256 signer
	signer := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "admin",
		"exp": time.Now().Add(time.Minute * 20).Unix(),
		"CustomUserInfo": struct {
			Name string
			Role string
		}{user.Username, "Member"}})

	//set claims

	tokenString, err := signer.SignedString(signKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		log.Printf("Error signing token: %v\n", err)
	}
	//create a token instance using the token string
	//response := Token{tokenString}
	cookie := &http.Cookie{
		Name:  "session",
		Value: tokenString,
		Path:  "/",
		//		UNCOMMENT WHEN DEPLOYED:
		//		Secure: true,
		//		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	//JsonResponse(response, w)

}

//ValidateTokenMiddleware is a AUTH TOKEN VALIDATION
func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//validate token
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err == nil {

		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorised access to this resource")
	}

}

//JSONResponse is a HELPER FUNCTIONS
func JSONResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func startServer() {
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)

	//PROTECTED ENDPOINTS
	/*	http.Handle("/resource/", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(ProtectedHandler)),
	))*/

	http.Handle("/api/getallvisits/", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(restHandler)),
	))
	http.Handle("/api/postavisit/", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(restHandler)),
	))
	http.Handle("/api/deleteavisit/", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(restHandler)),
	))

	log.Println("Now listening...")

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
