package model

type Friend struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	Group   string `json:"group"`
	Remark  string `json:"remark"`
	Created int64  `json:"created"`
}
