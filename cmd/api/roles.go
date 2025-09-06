package main

import (
	"net/http"
	"strconv"

	"github.com/faizisyellow/indocoffee/internal/repository"
	"github.com/faizisyellow/indocoffee/internal/service"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"github.com/go-chi/chi/v5"
)

// @Summary		Add user role
// @Description	Create new user role
// @Tags			Roles
// @Accept			json
// @Produce		json
// @Param			payload	body		service.CreateRoleRequest	true	"Payload create new role"
// @Success		201		{object}	main.Envelope{data=string,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		409		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/roles [post]
func (app *Application) CreateRolesHandler(w http.ResponseWriter, r *http.Request) {
	var req service.CreateRoleRequest
	if err := ReadHttpJson(w, r, &req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	res, err := app.Services.RolesService.Create(r.Context(), req.Serialize())
	if err != nil {
		errorValue := errorService.GetError(err)
		switch errorValue.E {
		case service.ErrConflictRole:
			ResponseClientError(w, r, err, http.StatusConflict)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, res, http.StatusCreated)
}

// @Summary		Get user roles
// @Description	Get All user roles
// @Tags			Roles
// @Produce		json
// @Success		200	{object}	main.Envelope{data=[]service.ResponseRolesFindAll,error=nil}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/roles [get]
func (app *Application) GetAllRolesHandler(w http.ResponseWriter, r *http.Request) {

	roles, err := app.Services.RolesService.FindAll(r.Context())
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
	}
	ResponseSuccess(w, r, roles, http.StatusOK)
}

// @Summary		Get user role
// @Description	Get user role by Id
// @Tags			Roles
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Role id"
// @Success		200	{object}	main.Envelope{data=service.ResponseRolesById,error=nil}
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/roles/{id} [get]
func (app *Application) GetRolesHandler(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	role, err := app.Services.RolesService.FindById(r.Context(), id)
	if err != nil {
		errorValue := errorService.GetError(err)
		switch errorValue.E {
		case service.ErrNotFoundRole:
			ResponseClientError(w, r, err, http.StatusNotFound)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}

		return
	}

	ResponseSuccess(w, r, role, http.StatusOK)
}

// @Summary		Update user role
// @Description	Update user role by Id
// @Tags			Roles
// @Accept			json
// @Produce		json
// @Param			id		path		int							true	"Role id"
// @Param			payload	body		service.RequestUpdateRole	true	"Role payload"
// @Success		200		{object}	main.Envelope{data=service.ResponseRolesById,error=nil}
// @Failure		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		404		{object}	main.Envelope{data=nil,error=string}
// @Failure		409		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/roles/{id} [patch]
func (app *Application) UpdateRolesHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	roleReq, err := app.Services.RolesService.FindById(ctx, id)
	if err != nil {
		errorValue := errorService.GetError(err)
		switch errorValue.E {
		case service.ErrNotFoundRole:
			ResponseClientError(w, r, err, http.StatusNotFound)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}

		return
	}

	var req service.RequestUpdateRole
	if err := ReadHttpJson(w, r, &req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(req); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	var role repository.RolesModel
	role.Id = roleReq.Id
	role.Name = roleReq.Name
	role.Level = roleReq.Level

	role = service.UpdateRolePayload(req.Serialize(), role)
	err = app.Services.RolesService.Update(ctx, id, role)
	if err != nil {
		errorValue := errorService.GetError(err)
		switch errorValue.E {
		case service.ErrConflictRole:
			ResponseClientError(w, r, err, http.StatusConflict)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, "success update role", http.StatusOK)
}

// @Summary		Delete user role
// @Description	Delete user role by Id
// @Tags			Roles
// @Accept			json
// @Produce		json
// @Param			id	path	int	true	"Role id"
// @Success		204
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/roles/{id} [delete]
func (app *Application) DeleteRolesHandler(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	err = app.Services.RolesService.Delete(r.Context(), id)
	if err != nil {
		errorValue := errorService.GetError(err)
		switch errorValue.E {
		case service.ErrNotFoundRole:
			ResponseClientError(w, r, err, http.StatusNotFound)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, nil, http.StatusNoContent)
}

// @Summary		Delete user roles
// @Description	Delete all user roles permanently
// @Tags			Roles
// @Success		204
// @Failure		404		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/roles [delete]
func (app *Application) TrashRolesHandler(w http.ResponseWriter, r *http.Request) {

}
