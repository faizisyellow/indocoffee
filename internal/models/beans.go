package models

type BeansModel struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsDelete bool   `json:"is_delete"`
}
