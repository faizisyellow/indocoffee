package dto

type CreateBeanRequest struct {
	Name string `json:"name" validate:"required,min=3,max=18"`
}

type UpdateBeanRequest struct {
	Name string `json:"name" validate:"required,min=3,max=18"`
}

type BeanResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
