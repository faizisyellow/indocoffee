package dto

type CreateFormRequest struct {
	Name string `json:"name" validate:"required,min=4,max=18"`
}

type FormResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UpdateFormRequest struct {
	Name string `json:"name" validate:"required,min=4,max=18"`
}
