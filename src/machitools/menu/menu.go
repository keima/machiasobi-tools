package menu

// +qbg
type MenuItem struct {
	Id          int64 `json:"id" datastore:"-" goon:"id"`
	Name        string `json:"name"`
	IconId      string `json:"icon"`
	State       string `json:"state"`
	Description string `json:"description" datastore:",noindex"`
	OrderIndex  int    `json:"order"`
	Enabled     bool   `json:"enabled"`
}

type MenuList []MenuItem
