package news

import (
	"github.com/knightso/base/gae/ds"
	"google.golang.org/appengine/datastore"
	"golang.org/x/net/context"
)

const (
	KindName = "News"
)

// ニュースのモデル
type NewsItem struct {
	ds.Meta
	Id       string `json:"id" datastore:"-"`
	Author   string `json:"-"`
	Title    string
	Article  string `datastore:",noindex"`
	IsPublic bool
}

func (item *NewsItem) Save(c context.Context, keyName string) error {
	key := datastore.NewKey(c, KindName, keyName, 0, nil)
	item.SetKey(key)
	return ds.Put(c, item)
}

func (item *NewsItem) Load(c context.Context, keyName string) error {
	key := datastore.NewKey(c, KindName, keyName, 0, nil)

	if err := ds.Get(c, key, item); err != nil {
		return err
	}

	// keyNameでもいいけれど。。。
	item.Id = item.GetKey().StringID()

	return nil
}

func LoadAll(c context.Context, first int, size int, publicOnly bool) (*[]NewsItem, error) {
	items := make([]NewsItem, 0, size)
	q := datastore.NewQuery(KindName).Order("-CreatedAt").Offset(first).Limit(size)

	if publicOnly {
		q = q.Filter("IsPublic = ", true)
	}

	if err := ds.ExecuteQuery(c, q, &items); err != nil {
		return nil, err
	}

	for index, item := range items {
		items[index].Id = item.GetKey().StringID()
	}

	return &items, nil
}
