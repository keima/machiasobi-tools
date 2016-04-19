package favorite

import (
	"errors"
	"machitools/customer"
	"github.com/knightso/base/gae/ds"
	"google.golang.org/appengine/datastore"
	"golang.org/x/net/context"
)

const KindName = "Fav"

type FavItem struct {
	ds.Meta
	CalendarID string `json:"calendarId"`
	EventID    string `json:"eventId"`
}

func (item *FavItem) Save(c context.Context, parent *datastore.Key) error {
	if parent.Kind() != customer.KindName {
		return errors.New("parent's kind name mismatch.")
	}
	if item.EventID == "" {
		return errors.New("Disallowed empty value `EventID`")
	}

	key := datastore.NewKey(c, KindName, item.EventID, 0, parent)
	item.SetKey(key)
	return ds.Put(c, item)
}

func (item *FavItem) Delete(c context.Context, parent *datastore.Key) error {
	if parent.Kind() != customer.KindName {
		return errors.New("parent's kind name mismatch.")
	}
	if item.EventID == "" {
		return errors.New("Disallowed empty value `EventID`")
	}

	key := datastore.NewKey(c, KindName, item.EventID, 0, parent)
	return ds.Delete(c, key)
}

func LoadAll(c context.Context, parent *datastore.Key) (*[]FavItem, error) {
	if parent.Kind() != customer.KindName {
		return nil, errors.New("parent's kind name mismatch.")
	}

	items := []FavItem{}
	q := datastore.NewQuery(KindName).Ancestor(parent).Order("-UpdatedAt")

	if err := ds.ExecuteQuery(c, q, &items); err != nil {
		return nil, err
	}

	return &items, nil
}
