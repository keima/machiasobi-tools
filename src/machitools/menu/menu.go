package menu

// +qbg
type MenuItem struct {
	Id          string `json:"id" datastore:"-" goon:"id"`
	Name        string `json:"string"`
	IconId      string `json:"icon"`
	State       string `json:"state"`
	Description string `json:"description"`
	OrderIndex  int    `json:"order"`
	Enabled     bool   `json:"enabled"`
}

type MenuList []MenuItem
