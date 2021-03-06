package customer

import (
	"errors"
	"github.com/knightso/base/gae/ds"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	"golang.org/x/net/context"
)

// Kind name of datastore
const KindName = "Customer"

type CustomerItem struct {
	ds.Meta
	Email string `json:"email"`
	ID    string `json:"userId"`
}

func (item *CustomerItem) Init(u *user.User) {
	item.Email = "" // u.Email is not correct by now...
	item.ID = u.ID
}

func (item *CustomerItem) Save(c context.Context) error {
	if item.ID == "" {
		return errors.New("ID is empty")
	}

	key := datastore.NewKey(c, KindName, item.ID, 0, nil)
	item.SetKey(key)
	return ds.Put(c, item)
}

func (item *CustomerItem) Load(c context.Context) error {
	key := datastore.NewKey(c, KindName, item.ID, 0, nil)
	if err := ds.Get(c, key, item); err != nil {
		return err
	}
	return nil
}
