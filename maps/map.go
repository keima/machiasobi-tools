package maps

import (
	"github.com/knightso/base/gae/model"

	"appengine"
	"appengine/datastore"
)

const kindNameMap = "Map"

type Map struct {
	model.Meta
	Id       string    `json:"id" datastore:"-"`
	Name     string    `json:"name"`
	IsPublic bool      `json:"isPublic"`
	MapItems []MapItem `json:"markers,omitempty" datastore:"-"`
}

func (item *Map) Save(c appengine.Context, keyName string) error {
	key := datastore.NewKey(c, kindNameMap, keyName, 0, nil)
	item.SetKey(key)
	return model.Put(c, item)
}

func (item *Map) Load(c appengine.Context, keyName string) error {
	key := datastore.NewKey(c, kindNameMap, keyName, 0, nil)

	if err := model.Get(c, key, item); err != nil {
		return err
	}

	item.Id = item.GetKey().StringID()

	return nil
}

func LoadAll(c appengine.Context, first, size int, publicOnly bool) (*[]Map, error) {
	items := make([]Map, 0, size)
	q := datastore.NewQuery(kindNameMap).Order("-UpdatedAt").Offset(first).Limit(size)

	if publicOnly {
		q = q.Filter("IsPublic = ", true)
	}

	if err := model.ExecuteQuery(c, q, &items); err != nil {
		return nil, err
	}

	for i, item := range items {
		items[i].Id = item.GetKey().StringID()
	}

	return &items, nil
}
