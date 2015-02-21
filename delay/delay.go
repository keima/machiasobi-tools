package delay

import (
	"errors"
	"strconv"

	"github.com/knightso/base/gae/ds"

	"appengine"
	"appengine/datastore"
)

type DelayItem struct {
	ds.Meta
	Id          string `datastore:"-" json:"id"`
	Delay       int    `json:"delay"`
	Message     string `json:"message"`
	Author      string `json:"-"`
	IsPostponed bool   `json:"isPostponed"`
}

var ErrItemNotFound = errors.New("Item is not found")

func (item *DelayItem) Save(c appengine.Context, placeName string) error {
	key := datastore.NewIncompleteKey(c, kindName(placeName), nil)
	item.SetKey(key)
	return ds.Put(c, item)
}

func LoadLatest(c appengine.Context, placeName string) (*DelayItem, error) {
	items := make([]DelayItem, 0, 1)

	q := datastore.NewQuery(kindName(placeName))
	q = q.Order("-CreatedAt").Limit(1)

	if err := ds.ExecuteQuery(c, q, &items); err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, ErrItemNotFound
	}

	for i, item := range items {
		items[i].Id = strconv.FormatInt(item.GetKey().IntID(), 10)
	}

	return &items[0], nil
}

const kindNamePrefix = "Delay-"

func kindName(placeName string) string {
	return kindNamePrefix + placeName
}
