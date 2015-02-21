package maps

import (
	"strconv"

	"github.com/knightso/base/gae/model"

	"appengine"
	"appengine/datastore"
)

const kindNameMapItem = "MapItem"

type MapItem struct {
	model.Meta
	Id          string       `json:"id" datastore:"-"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Coords      GeoPointJson `json:"coords"`
	Order       int          `json:"order"`
}

func (item *MapItem) Save(c appengine.Context, parent *datastore.Key) error {
	key := datastore.NewIncompleteKey(c, kindNameMapItem, parent)
	item.SetKey(key)
	return model.Put(c, item)
}

func (item *MapItem) Update(c appengine.Context, parent *datastore.Key, keyId int64) error {
	key := datastore.NewKey(c, kindNameMapItem, "", keyId, parent)
	item.SetKey(key)
	return model.Put(c, item)
}

func (item *MapItem) Load(c appengine.Context, parent *datastore.Key, keyName string) error {
	keyId, err := strconv.ParseInt(keyName, 10, 64)
	if err != nil {
		return err
	}

	key := datastore.NewKey(c, kindNameMapItem, "", keyId, parent)
	if err := model.Get(c, key, item); err != nil {
		return err
	}

	item.Id = strconv.FormatInt(item.GetKey().IntID(), 10)

	return nil
}

func LoadAllMapItem(c appengine.Context, parent *datastore.Key) (*[]MapItem, error) {
	items := make([]MapItem, 0, 20) // TODO: magic number...
	q := datastore.NewQuery(kindNameMapItem).Ancestor(parent)
	q = q.Order("CreatedAt")

	if err := model.ExecuteQuery(c, q, &items); err != nil {
		return nil, err
	}

	for i, item := range items {
		items[i].Id = strconv.FormatInt(item.GetKey().IntID(), 10)
	}

	return &items, nil
}
