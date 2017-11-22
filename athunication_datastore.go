package hellodatastore

import (
	"encoding/json"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

// USERENTITY User Is using for store the users information
const USERENTITY = "User"

// NameSpace -kashyak- Is the defalut Namespace
//const NameSpace = "-kashyak-"

//var tpl *template.Template

//UserEntiry is a struct
type UserEntiry struct {
	UserName  string `datastore:"-"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	PassWord  string `json:"password"`
}
type sessionData struct {
	UserEntiry
	LoggedIn  bool
	LoginFail bool
}

//UserDatasotre is to connect the datastore
type UserDatasotre struct {
	Ctx  context.Context
	Keys []*datastore.Key
	User *UserEntiry
}

//NewUserDatasotre is have the values for connecting datastore
func NewUserDatasotre(r *http.Request) *UserDatasotre {
	ed := new(UserDatasotre)
	ed.Ctx = appengine.NewContext(r)
	return ed
}

func (ed *UserDatasotre) build(uname string, fname string, lname string, email string, pass string) error {
	ed.User = &UserEntiry{
		UserName:  uname,
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		PassWord:  pass,
	}

	var err error
	//ed.Keys, err = getKeysByName(ed)
	//log.Printf("%#v ed q values", ed.Keys)
	return err

}

func (ed *UserDatasotre) login(ud *UserDatasotre, res http.ResponseWriter, req *http.Request) error {
	ctx := appengine.NewContext(req)
	log.Infof(ctx, "ed.User.UserName==>"+"%v", ed.User.UserName)
	key := datastore.NewKey(ctx, USERENTITY, ed.User.UserName, 0, nil)
	var user UserEntiry
	err := datastore.Get(ctx, key, &user)

	user.UserName = ed.User.UserName
	createSession(user, res, req)
	//http.Redirect(res, req, "/", http.StatusSeeOther)
	callmemcache(req)

	return err
}

//Put Method using to put the user values to datastore
func (ed *UserDatasotre) createUser(ud *UserDatasotre, w http.ResponseWriter, r *http.Request) error {

	ctx := appengine.NewContext(r)
	if len(ed.Keys) == 0 {
		keys := make([]*datastore.Key, 1)

		keys[0] = datastore.NewKey(ctx, USERENTITY, ed.User.UserName, 0, nil)
		//keys[0] = datastore.NewIncompleteKey(ctx, USERENTITY, nil)
		ed.Keys = keys

	}
	_, err := datastore.Put(ctx, ed.Keys[0], ed.User)
	//createSession(ed.User, w, r)
	return err
}

func createSession(user UserEntiry, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// SET COOKIE
	id := uuid.NewV4()
	cookie := &http.Cookie{
		Name:  "session",
		Value: id.String(),
		Path:  "/",
		//		UNCOMMENT WHEN DEPLOYED:
		//		Secure: true,
		//		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	// SET MEMCACHE session data (sd)
	json, err := json.Marshal(user)
	if err != nil {
		log.Errorf(ctx, "error marshalling during user creation: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	sd := memcache.Item{
		Key:   id.String(),
		Value: json,
		//		Expiration: time.Duration(20*time.Minute),
		Expiration: time.Duration(20 * time.Second),
	}
	memcache.Set(ctx, &sd)

}
func callmemcache(r *http.Request) {
	session := getSession(r)

	if len(session.Value) > 0 {
		var sd sessionData
		json.Unmarshal(session.Value, &sd)
		sd.LoggedIn = true

	}
}
