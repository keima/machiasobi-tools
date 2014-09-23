package news

import (
	"appengine"
	"appengine/datastore"
	"time"
)

const (
	KindName = "News"
)

// ニュースのモデル
type NewsItem struct {
	Author   string
	Title    string
	Article  string
	Date     time.Time
	IsPublic bool
}

func (news *NewsItem) save(c appengine.Context, keyName string) (*datastore.Key, error) {
	key := datastore.NewKey(c, KindName, keyName, 0, nil)
	return datastore.Put(c, key, news)
}

func (NewsItem) loadAll(c appengine.Context, first int, size int, publicOnly bool) (*[]NewsItem, error) {
	items := make([]NewsItem, 0, size)
	query := datastore.NewQuery(KindName).Order("-Date").Offset(first).Limit(size)

	if (publicOnly) {
		query = query.Filter("IsPublic = ", true)
	}

	if _, err := query.GetAll(c, &items); err != nil {
		return nil, err
	}

	return &items, nil
}

func (NewsItem) load(c appengine.Context, keyName string) (*NewsItem, error) {
	item := NewsItem{}
	key := datastore.NewKey(c, KindName, keyName, 0, nil)

	if err := datastore.Get(c, key, &item); err != nil {
		return nil, err
	}
	return &item, nil
}
