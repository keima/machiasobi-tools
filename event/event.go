package event

import (
	"log"
	"strconv"
	"time"

	"github.com/knightso/base/gae/ds"

	"appengine"
	"appengine/datastore"
)

type EventItem struct {
	ds.Meta
	Id      string    `datastore:"-"`
	Title   string    `json:"title"`
	Place   string    `json:"place"`
	Message string    `json:"message" datastore:",noindex"`
	StartAt time.Time `json:"startAt"`

	Author string `json:"-"`

	IsPublic   bool `json:"isPublic"`
	IsRunning  bool `json:"isRunning"`
	IsFinished bool `json:"isFinished"`
}

const kindName = "Event"

func (item *EventItem) Save(c appengine.Context) error {
	key := datastore.NewIncompleteKey(c, kindName, nil)
	item.SetKey(key)
	return ds.Put(c, item)
}

func (item *EventItem) SaveUpdate(c appengine.Context, intID int64) error {
	key := datastore.NewKey(c, kindName, "", intID, nil)
	item.SetKey(key)
	return ds.Put(c, item)
}

func (item *EventItem) Load(c appengine.Context, keyName string) error {
	keyId, err := strconv.ParseInt(keyName, 10, 64)
	if err != nil {
		return err
	}

	key := datastore.NewKey(c, kindName, "", keyId, nil)

	if err := ds.Get(c, key, item); err != nil {
		return err
	}

	item.Id = strconv.FormatInt(item.GetKey().IntID(), 10)

	return nil
}

func LoadAll(c appengine.Context, first int, size int, rangeStart time.Time, rangeEnd time.Time, publicOnly bool) (*[]EventItem, error) {
	items := make([]EventItem, 0, size)
	q := datastore.NewQuery(kindName)

	log.Println("hoge")

	if publicOnly {
		q = q.Filter("IsPublic = ", true)
		log.Println("p")
	}

	if !rangeStart.IsZero() || !rangeEnd.IsZero() {
		q = q.Order("-StartAt")

		if !rangeStart.IsZero() {
			q = q.Filter("StartAt >", rangeStart)
			log.Println("s")
		}

		if !rangeEnd.IsZero() {
			q = q.Filter("StartAt <", rangeEnd)
			log.Println("e")
		}
	}

	q = q.Order("-UpdatedAt").Offset(first).Limit(size)

	if err := ds.ExecuteQuery(c, q, &items); err != nil {
		return nil, err
	}

	for index, item := range items {
		items[index].Id = strconv.FormatInt(item.GetKey().IntID(), 10)
	}

	return &items, nil
}
