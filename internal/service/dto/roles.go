package dto

type CreateRoleRequest struct {
	Name  string `json:"name" validate:"required,min=4,max=16"`
	Level int    `json:"level" validate:"required,min=1,max=5"`
}

type RolesResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}

type UpdateRoleRequest struct {
	Name  *string `json:"name" validate:"omitempty,min=4,max=16"`
	Level *int    `json:"level" validate:"omitempty,min=1,max=5"`
}
