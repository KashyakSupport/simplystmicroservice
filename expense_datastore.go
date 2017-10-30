package hellodatastore

import (
	"log"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

const ENTITYNAME = "User"
const NameSpace = "-kashyak-"

type ExpenseEntiry struct {
	Name string `json:"string"`
	//Price    int    `json:"price"`
	LastName string `json:"lastname"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
	Email    string `json:"email"`
}

type ExpenseDatasotre struct {
	Ctx     context.Context
	Keys    []*datastore.Key
	Expense *ExpenseEntiry
}

func NewExpenseDatasotre(r *http.Request) *ExpenseDatasotre {
	ed := new(ExpenseDatasotre)
	ed.Ctx = appengine.NewContext(r)
	//ctx, _ := appengine.Namespace(ed.Ctx, "-global-")
	return ed
}

func (ed *ExpenseDatasotre) build(name string, lname string, uname string, pass string, email string) error {
	ed.Expense = &ExpenseEntiry{
		Name: name,
		//Price:    price,
		LastName: lname,
		UserName: uname,
		PassWord: pass,
		Email:    email,
	}

	log.Printf("%#v", ed.Expense)
	var err error
	/*ed.Keys, err = getKeysByName(ed)
	return err*/
	return err

}

/*func (ed *ExpenseDatasotre) Get() error {
	if len(ed.Keys) == 0 {
		return errors.New("Item to be acquired does not exist")
	}
	return datastore.Get(ed.Ctx, ed.Keys[0], ed.Expense)
}*/

func (ed *ExpenseDatasotre) Put() error {
	var ctx context.Context
	log.Printf("%#v", ed.Keys)
	log.Printf("%#v", len(ed.Keys))
	if len(ed.Keys) == 0 {
		keys := make([]*datastore.Key, 1)
		log.Printf("%#v", keys)
		log.Printf("%#v", ed.Ctx)
		log.Printf("%#v", ENTITYNAME)
		ctx, _ = appengine.Namespace(ed.Ctx, NameSpace)
		keys[0] = datastore.NewIncompleteKey(ctx, ENTITYNAME, nil)
		ed.Keys = keys
		log.Printf("%#v", keys)
		log.Printf("%#v", ed.Keys)
	}

	_, err := datastore.Put(ctx, ed.Keys[0], ed.Expense)
	return err
}

/*func (ed *ExpenseDatasotre) Delete() error {
	if len(ed.Keys) == 0 {
		return errors.New("does not exist item to be deleted")
	}

	option := &datastore.TransactionOptions{XG: true}
	return datastore.RunInTransaction(ed.Ctx, func(c context.Context) error {
		return datastore.DeleteMulti(c, ed.Keys)
	}, option)
}*/

/*func getKeysByName(ed *ExpenseDatasotre) ([]*datastore.Key, error) {
	q := datastore.NewQuery(ENTITYNAME).Filter("Name =", ed.Expense.Name)

	var expenses []ExpenseEntiry
	return q.GetAll(ed.Ctx, &expenses)

}*/
