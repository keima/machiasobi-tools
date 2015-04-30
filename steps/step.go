package steps

import (
	"github.com/keima/base/gae/ds"

	"appengine"
	"appengine/datastore"
)

type StepItem struct {
	ds.Meta
	Id          int64  `json:"id" datastore:"-"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	ShowTitle   bool   `json:"showTitle"`
	Description string `json:"description" datastore:",noindex"`
	Path        string `json:"path"` // ← not use
	PartialId   string `json:"partialId"`
	Author      string `json:"-"`
	IsPublic    bool   `json:"isPublic"`
	Order       int    `json:"order"`
}

const (
	KindName = "Step"
)

var AllowedType = [...]string{"partial", "html", "markdown"}

func (item *StepItem) Save(c appengine.Context) error {
	key := datastore.NewIncompleteKey(c, KindName, nil)
	item.SetKey(key)
	return ds.Put(c, item)
}

func (new *StepItem) Update(c appengine.Context, keyId int64) error {
	old := StepItem{}
	if err := old.Load(c, keyId); err != nil {
		return err
	}

	old.Id = new.Id
	old.Type = new.Type
	old.Title = new.Title
	old.ShowTitle = new.ShowTitle
	old.Description = new.Description
	old.Path = new.Path
	old.PartialId = new.PartialId
	old.Author = new.Author
	old.IsPublic = new.IsPublic
	old.Order = new.Order

	key := datastore.NewKey(c, KindName, "", keyId, nil)
	old.SetKey(key)
	return ds.Put(c, &old)
}

func (item *StepItem) Load(c appengine.Context, keyId int64) error {
	key := datastore.NewKey(c, KindName, "", keyId, nil)

	if err := ds.Get(c, key, item); err != nil {
		return err
	}

	// item.Id = strconv.FormatInt(item.GetKey().IntID(), 10)
	item.Id = item.GetKey().IntID()

	return nil
}

func LoadAll(c appengine.Context, first, size int, private bool) (*[]StepItem, error) {
	items := make([]StepItem, 0, size)
	q := datastore.NewQuery(KindName).Order("Order").Offset(first).Limit(size)

	if !private {
		q = q.Filter("IsPublic = ", true)
	}

	if err := ds.ExecuteQuery(c, q, &items); err != nil {
		return nil, err
	}

	for i, item := range items {
		// items[i].Id = strconv.FormatInt(item.GetKey().IntID(), 10)
		items[i].Id = item.GetKey().IntID()
	}

	return &items, nil
}
