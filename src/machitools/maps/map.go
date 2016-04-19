package maps

import (
	"github.com/knightso/base/gae/ds"
	"google.golang.org/appengine/datastore"
	"golang.org/x/net/context"
)

const kindNameMap = "Map"

type Map struct {
	ds.Meta
	Id       string    `json:"id" datastore:"-"`
	Name     string    `json:"name"`
	IsPublic bool      `json:"isPublic"`
	MapItems []MapItem `json:"markers,omitempty" datastore:"-"`
}

func (item *Map) Save(c context.Context, keyName string) error {
	key := datastore.NewKey(c, kindNameMap, keyName, 0, nil)
	item.SetKey(key)
	return ds.Put(c, item)
}

func (item *Map) Load(c context.Context, keyName string) error {
	key := datastore.NewKey(c, kindNameMap, keyName, 0, nil)

	if err := ds.Get(c, key, item); err != nil {
		return err
	}

	item.Id = item.GetKey().StringID()

	return nil
}

func LoadAll(c context.Context, first, size int, private bool) (*[]Map, error) {
	items := make([]Map, 0, size)
	q := datastore.NewQuery(kindNameMap).Order("-UpdatedAt").Offset(first).Limit(size)

	if !private {
		q = q.Filter("IsPublic = ", true)
	}

	if err := ds.ExecuteQuery(c, q, &items); err != nil {
		return nil, err
	}

	for i, item := range items {
		items[i].Id = item.GetKey().StringID()
	}

	return &items, nil
}
