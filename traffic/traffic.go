package traffic

import (
	"appengine"
	"appengine/datastore"
	"time"
	"errors"
)

const KindNamePrefix = "Traffic-"

// Trafficのモデル
type TrafficItem struct {
	Waiting int
	Message string
	Author  string
	Direction int
	Date    time.Time
}

// Directionの方向
const (
	DirectionError = -1

	DirectionInbound = 0
	DirectionOutbound = 1
)

// 見つからなかったよエラー
var ErrItemNotFound = errors.New("Item is not found")



func (item *TrafficItem) save(c appengine.Context, trafficName string) (*datastore.Key, error) {
	key := datastore.NewIncompleteKey(c, KindName(trafficName), nil)
	return datastore.Put(c, key, item)
}

func (TrafficItem) loadLatest(c appengine.Context, trafficName string, direction int) (*TrafficItem, error) {
	items := make([]TrafficItem, 0, 1)

	q := datastore.NewQuery(KindName(trafficName))
	q = q.Filter("Direction = ", direction).Order("-Date").Limit(1)

	if _,err := q.GetAll(c, &items); err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil ,ErrItemNotFound
	}

	return &items[0], nil
}

func KindName(trafficName string) string {
	return KindNamePrefix + trafficName
}
