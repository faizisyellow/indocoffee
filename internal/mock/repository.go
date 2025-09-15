package mock

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/faizisyellow/indocoffee/internal/repository"
	"github.com/faizisyellow/indocoffee/internal/utils"
)

type UsersRepositoryMock struct {
}

func (u *UsersRepositoryMock) Insert(ctx context.Context, tx *sql.Tx, usr repository.UserModel) (int, error) {
	existingUser := repository.UserModel{
		Email: "lizzymcalpine@test.com",
	}

	if usr.Email == existingUser.Email {
		return 0, errors.New("account already exist")
	}

	return usr.Id, nil
}

func (u *UsersRepositoryMock) GetById(ctx context.Context, id int) (repository.UserModel, error) {

	existingUser := repository.UserModel{
		Id:       21,
		Username: "lizzy",
		Email:    "lizzymcalpine@test.com",
	}

	if existingUser.Id != id {
		return repository.UserModel{}, errors.New("user not found")
	}

	return existingUser, nil
}

func (u *UsersRepositoryMock) GetByEmail(ctx context.Context, email string) (repository.UserModel, error) {
	existingUser := repository.UserModel{
		Id:        21,
		Email:     "lizzymcalpine@test.com",
		IsActive:  utils.BoolToPoint(true),
		CreatedAt: time.Time{},
	}
	existingUser.Password.ParseFromPassword("HelloWorld$123")

	userNotAcivated := repository.UserModel{
		Id:        29,
		Email:     "coolemailname@test.com",
		IsActive:  utils.BoolToPoint(false),
		CreatedAt: time.Time{},
	}
	userNotAcivated.Password.ParseFromPassword("HelloWorld$123")

	if userNotAcivated.Email == email {
		return userNotAcivated, nil
	}

	if existingUser.Email != email && userNotAcivated.Email != email {
		return repository.UserModel{}, errors.New("user not found")
	}

	return existingUser, nil
}

func (u *UsersRepositoryMock) Update(ctx context.Context, tx *sql.Tx, usr repository.UserModel) error {

	return nil
}

func (u *UsersRepositoryMock) Delete(ctx context.Context, tx *sql.Tx, id int) error {

	return nil
}

type InvitationRepositoryMock struct {
}

func (i *InvitationRepositoryMock) Insert(ctx context.Context, tx *sql.Tx, invt repository.InvitationModel) error {

	return nil
}
func (i *InvitationRepositoryMock) Get(ctx context.Context, tx *sql.Tx, token string) (int, error) {

	if token == "" {
		return 0, errors.New("token invitation not found")
	}

	// user's id
	return 21, nil
}
func (i *InvitationRepositoryMock) DeleteByUserId(ctx context.Context, tx *sql.Tx, usrid int) error {

	return nil
}

type BeansRepositoryMock struct{}

func (b *BeansRepositoryMock) Insert(ctx context.Context, nw repository.BeansModel) error {
	existBean := repository.BeansModel{
		Name: "robusta",
	}

	if existBean.Name == nw.Name {
		conflictErrorMsg := strings.Builder{}
		conflictErrorMsg.WriteString("Error 1062 (23000): ")
		conflictErrorMsg.WriteString("can not duplicat row")

		return errors.New(conflictErrorMsg.String())
	}
	return nil
}
func (b *BeansRepositoryMock) GetAll(ctx context.Context) ([]repository.BeansModel, error) {

	beans := []repository.BeansModel{
		{
			Id:   1,
			Name: "arabica",
		},
	}

	return beans, nil
}
func (b *BeansRepositoryMock) GetById(ctx context.Context, id int) (repository.BeansModel, error) {

	foundBean := repository.BeansModel{
		Id:   1,
		Name: "arabica",
	}

	if id != foundBean.Id {
		return repository.BeansModel{}, sql.ErrNoRows
	}

	return foundBean, nil
}
func (b *BeansRepositoryMock) Update(ctx context.Context, nw repository.BeansModel) error {

	// let's say in db already have
	// this bean with this name
	if nw.Name == "robusta" {
		conflictErrorMsg := strings.Builder{}
		conflictErrorMsg.WriteString("Error 1062 (23000): ")
		conflictErrorMsg.WriteString("can not duplicat row")
		return errors.New(conflictErrorMsg.String())
	}

	return nil
}

func (b *BeansRepositoryMock) Delete(ctx context.Context, id int) error {
	return nil
}

func (b *BeansRepositoryMock) DestroyMany(ctx context.Context) error {
	return nil
}

type FormsRepositoryMock struct{}

func (f *FormsRepositoryMock) Insert(ctx context.Context, nw repository.FormsModel) error {
	exstForm := repository.FormsModel{
		Name: "grounded",
	}

	if exstForm.Name == nw.Name {
		conflictErrorMsg := strings.Builder{}
		conflictErrorMsg.WriteString("Error 1062 (23000): ")
		conflictErrorMsg.WriteString("can not duplicate row")
		return errors.New(conflictErrorMsg.String())
	}

	return nil
}
func (f *FormsRepositoryMock) GetAll(ctx context.Context) ([]repository.FormsModel, error) {

	exstForms := []repository.FormsModel{
		{
			Id:   1,
			Name: "grounded",
		},
	}

	return exstForms, nil
}
func (f *FormsRepositoryMock) GetById(ctx context.Context, id int) (repository.FormsModel, error) {

	exstForm := repository.FormsModel{
		Id:   1,
		Name: "grounded",
	}

	if exstForm.Id != id {
		return repository.FormsModel{}, sql.ErrNoRows
	}

	return exstForm, nil
}
func (f *FormsRepositoryMock) Update(ctx context.Context, nw repository.FormsModel) error {

	exstForm := repository.FormsModel{
		Name: "grounded",
	}

	if exstForm.Name == nw.Name {
		conflictErrorMsg := strings.Builder{}
		conflictErrorMsg.WriteString("Error 1062 (23000): ")
		conflictErrorMsg.WriteString("can not duplicate row")
		return errors.New(conflictErrorMsg.String())
	}

	return nil
}
func (f *FormsRepositoryMock) Delete(ctx context.Context, id int) error {

	return nil
}
func (f *FormsRepositoryMock) DestroyMany(ctx context.Context) error {

	return nil
}

type RolesRepositoryMock struct {
}

func (r *RolesRepositoryMock) Insert(ctx context.Context, nw repository.RolesModel) error {
	extRole := repository.RolesModel{
		Name: "admin",
	}

	if extRole.Name == nw.Name {
		conflictErrorMsg := strings.Builder{}
		conflictErrorMsg.WriteString("Error 1062 (23000): ")
		conflictErrorMsg.WriteString("can not duplicate row")
		return errors.New(conflictErrorMsg.String())
	}
	return nil
}

func (r *RolesRepositoryMock) GetAll(ctx context.Context) ([]repository.RolesModel, error) {

	extRole := []repository.RolesModel{
		{
			Id:    1,
			Name:  "admin",
			Level: 3,
		},
	}

	return extRole, nil
}

func (r *RolesRepositoryMock) GetById(ctx context.Context, id int) (repository.RolesModel, error) {
	extRole := repository.RolesModel{
		Id:    1,
		Name:  "admin",
		Level: 3,
	}

	if extRole.Id != id {
		return repository.RolesModel{}, sql.ErrNoRows
	}

	return extRole, nil
}

func (r *RolesRepositoryMock) DestroyMany(ctx context.Context) error {

	return nil
}

func (r *RolesRepositoryMock) Update(ctx context.Context, nw repository.RolesModel) error {

	extRole := repository.RolesModel{
		Name: "manager",
	}

	if extRole.Name == nw.Name {
		conflictErrorMsg := strings.Builder{}
		conflictErrorMsg.WriteString("Error 1062 (23000): ")
		conflictErrorMsg.WriteString("can not duplicate row")
		return errors.New(conflictErrorMsg.String())
	}

	return nil
}

func (r *RolesRepositoryMock) Delete(ctx context.Context, id int) error {
	return nil
}
