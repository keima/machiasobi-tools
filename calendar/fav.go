package calendar

import (
	"errors"

	"github.com/keima/machiasobi-tools/customer"
	"github.com/keima/base/gae/ds"

	"appengine"
	"appengine/datastore"
)

const KindName = "Fav"

type FavItem struct {
	ds.Meta
	CalendarID string `json:"calendarId"`
	EventID    string `json:"eventId"`
}

func (item *FavItem) Save(c appengine.Context, parent *datastore.Key) error {
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

func (item *FavItem) Delete(c appengine.Context, parent *datastore.Key) error {
	if parent.Kind() != customer.KindName {
		return errors.New("parent's kind name mismatch.")
	}
	if item.EventID == "" {
		return errors.New("Disallowed empty value `EventID`")
	}

	key := datastore.NewKey(c, KindName, item.EventID, 0, parent)
	return ds.Delete(c, key)
}

func LoadAll(c appengine.Context, parent *datastore.Key) (*[]FavItem, error) {
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
