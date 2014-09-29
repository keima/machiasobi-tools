package news

import (
	"appengine"
	"appengine/datastore"
	"github.com/knightso/base/gae/model"
)

const (
	KindName = "News"
)

// ニュースのモデル
type NewsItem struct {
	model.Meta
	Id       string `datastore:"-"`
	Author   string `json:"-"`
	Title    string
	Article  string `datastore:",noindex"`
	IsPublic bool
}

func (item *NewsItem) Save(c appengine.Context, keyName string) error {
	key := datastore.NewKey(c, KindName, keyName, 0, nil)
	item.SetKey(key)
	return model.Put(c, item)
}

func (news *NewsItem) Load(c appengine.Context, keyName string) error {
	key := datastore.NewKey(c, KindName, keyName, 0, nil)

	if err := model.Get(c, key, news); err != nil {
		return err
	}

	// keyNameでもいいけれど。。。
	news.Id = news.GetKey().StringID()

	return nil
}

func LoadAll(c appengine.Context, first int, size int, publicOnly bool) (*[]NewsItem, error) {
	items := make([]NewsItem, 0, size)
	q := datastore.NewQuery(KindName).Order("-UpdatedAt").Offset(first).Limit(size)

	if (publicOnly) {
		q = q.Filter("IsPublic = ", true)
	}

	if err := model.ExecuteQuery(c, q, &items); err != nil {
		return nil, err
	}

	for index, item := range items {
		items[index].Id = item.GetKey().StringID();
	}

	return &items, nil
}
