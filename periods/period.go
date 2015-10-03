package periods

import (
	"github.com/mjibson/goon"
	"time"
	"net/http"

	"google.golang.org/appengine/datastore"
)

type PeriodItem struct {
	Id       int64 `datastore:"-" goon:"id" json:"id,omitempty"`
	Date     time.Time `json:"date"`
	IsActive bool `json:"isActive"`
}

type PeriodItemList []*PeriodItem

func (item *PeriodItem) Save(r *http.Request) error {
	g := goon.NewGoon(r)
	_, err := g.Put(item)
	return err
}

// 必ず item := PeriodItem{Id: 12345} してから使うこと
func (item *PeriodItem) Load(r *http.Request) error {
	g := goon.NewGoon(r)
	err := g.Get(item)
	return err
}

func (items *PeriodItemList) LoadActive(r *http.Request) error {
	g := goon.NewGoon(r)
	q := datastore.NewQuery(g.Kind(new(PeriodItem))).Order("Date").Filter("IsActive =", true)
	_, err := g.GetAll(q, items)
	return err
}