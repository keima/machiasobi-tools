package customer

import (
	"errors"

	"github.com/knightso/base/gae/ds"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

// Kind name of datastore
const KindName = "Customer"

type CustomerItem struct {
	ds.Meta
	Email string `json:"email"`
	ID    string `json:"userId"`
}

func (item *CustomerItem) Init(u *user.User) {
	item.Email = u.Email
	item.ID = u.ID
}

func (item *CustomerItem) Save(c appengine.Context) error {
	if item.ID == "" {
		return errors.New("ID is empty")
	}

	key := datastore.NewKey(c, KindName, item.ID, 0, nil)
	item.SetKey(key)
	return ds.Put(c, item)
}

func (item *CustomerItem) Load(c appengine.Context) error {
	key := datastore.NewKey(c, KindName, item.ID, 0, nil)
	if err := ds.Get(c, key, item); err != nil {
		return err
	}
	return nil
}
