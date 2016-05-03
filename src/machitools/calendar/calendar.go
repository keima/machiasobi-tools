package calendar

import (
	"net/http"
	"github.com/mjibson/goon"
	"google.golang.org/appengine/datastore"
)

// +qbg
type CalendarItem struct {
	Id              string `json:"id" datastore:"-" goon:"id"`
	Name            string `json:"name"`
	CalendarId      string `json:"calendarId"`
	IsSticky        bool   `json:"isSticky"`
	Order           int    `json:"order"`
	Enabled         bool   `json:"enabled,omitempty"`
}

type CalendarItemList []*CalendarItem

func (item *CalendarItem) Save(r *http.Request) error {
	g := goon.NewGoon(r)
	_, err := g.Put(item)
	return err
}

func (item *CalendarItem) Load(r *http.Request, id string) error {
	item.Id = id

	g := goon.NewGoon(r)
	err := g.Get(item)
	return err
}

func (item *CalendarItem) Delete(r *http.Request) error {
	g := goon.NewGoon(r)
	key := g.Key(item)
	return g.Delete(key)
}

func (items *CalendarItemList) LoadAll(r *http.Request) error {
	g := goon.NewGoon(r)
	q := query(g)
	_, err := g.GetAll(q, items)
	return err
}

func (items *CalendarItemList) LoadEnabled(r *http.Request) error {
	g := goon.NewGoon(r)
	q := query(g).Filter("Enabled =", true)
	_, err := g.GetAll(q, items)
	return err
}

func (CalendarItemList) Count(r *http.Request) (int,error) {
	g := goon.NewGoon(r)
	q := query(g)
	return g.Count(q)
}

func query(g *goon.Goon) *datastore.Query {
	return datastore.NewQuery(g.Kind(new(CalendarItem))).Order("Order")
}