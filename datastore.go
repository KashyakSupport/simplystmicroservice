package simplystmicroservice

import (
	"errors"
	"net/http"

	"google.golang.org/appengine"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

//ENTITYNAME is the Kind Name
const ENTITYNAME = "User"

// UserEntiry is the form values
type UserEntiry struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lasttname"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
}

//UserDatastore is to entry
type UserDatastore struct {
	Ctx  context.Context
	Keys []*datastore.Key
	User *UserEntiry
}

//NewUserDatastore is a new data base
func NewUserDatastore(r *http.Request) *UserDatastore {
	ud := new(UserDatastore)
	ud.Ctx = appengine.NewContext(r)
	return ud
}

func (ud *UserDatastore) build(firstname string, lastname string, username string, password string) error {

	ud.User = &UserEntiry{
		FirstName: firstname,
		LastName:  lastname,
		UserName:  username,
		Password:  password,
	}
	//log.Printf("%#v", ud.User)
	//var err error
	//ud.Keys,
	return nil
}

//Get vlrkvj
func (ud *UserDatastore) Get() error {
	if len(ud.Keys) == 0 {
		return errors.New("Hi eeror from get")
	}
	return datastore.Get(ud.Ctx, ud.Keys[0], ud.User)
}

//Put is method to insert values
func (ud *UserDatastore) Put() error {
	if len(ud.Keys) == 0 {
		keys := make([]*datastore.Key, 1)
		keys[0] = datastore.NewIncompleteKey(ud.Ctx, ENTITYNAME, nil)
		ud.Keys = keys
	}

	_, err := datastore.Put(ud.Ctx, ud.Keys[0], ud.User)
	return err
}
