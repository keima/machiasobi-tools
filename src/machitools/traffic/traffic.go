package traffic

import (
	"errors"
	"strconv"
	"github.com/knightso/base/gae/ds"
	"google.golang.org/appengine/datastore"
	"golang.org/x/net/context"
)

// Trafficのモデル
type TrafficItem struct {
	ds.Meta
	Id        string `datastore:"-"`
	Waiting   int
	Message   string
	Author    string `json:"-"`
	Direction int
}

// Directionの方向
const (
	DirectionError = -1

	DirectionInbound  = 0
	DirectionOutbound = 1
)

// 見つからなかったよエラー
var ErrItemNotFound = errors.New("Item is not found")

func (item *TrafficItem) Save(c context.Context, trafficName string) error {
	key := datastore.NewIncompleteKey(c, kindName(trafficName), nil)
	item.SetKey(key)
	return ds.Put(c, item)
}

func LoadLatest(c context.Context, trafficName string, direction int) (*TrafficItem, error) {
	items := make([]TrafficItem, 0, 1)

	q := datastore.NewQuery(kindName(trafficName))
	q = q.Filter("Direction = ", direction).Order("-CreatedAt").Limit(1)

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

const kindNamePrefix = "Traffic-"

func kindName(trafficName string) string {
	return kindNamePrefix + trafficName
}
