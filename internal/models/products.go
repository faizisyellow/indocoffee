package models

type Product struct {
	Id         int     `json:"id"`
	Roasted    string  `json:"roasted"`
	Price      float64 `json:"price"`
	Quantity   int     `json:"quantity"`
	Image      string  `json:"image"`
	BeanId     int     `json:"bean_id"`
	FormId     int     `json:"form_id"`
	BeansModel `json:"bean"`
	FormsModel `json:"form"`
}
