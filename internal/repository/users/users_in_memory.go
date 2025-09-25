package users

import (
	"context"
	"database/sql"
	"time"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/utils"
)

type InMemoryUsers struct {
	Users []models.User
}

func (u *InMemoryUsers) Insert(ctx context.Context, _ *sql.Tx, usr models.User) (int, error) {

	nextID := 1
	for _, user := range u.Users {
		if user.Id >= nextID {
			nextID = user.Id + 1
		}
	}

	newUser := models.User{
		Id:        nextID,
		Username:  usr.Username,
		Email:     usr.Email,
		Password:  usr.Password,
		IsActive:  utils.BoolToPoint(false),
		CreatedAt: time.Time{},
	}
	u.Users = append(u.Users, newUser)

	return newUser.Id, nil
}

func (u *InMemoryUsers) GetById(ctx context.Context, id int) (models.User, error) {
	var usr = models.User{}
	for _, user := range u.Users {
		if user.Id == id {
			usr = user
		}
	}

	return usr, nil
}

func (u *InMemoryUsers) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var usr = models.User{}
	for _, user := range u.Users {
		if user.Email == email {
			usr = user
		}
	}
	return usr, nil
}

func (u *InMemoryUsers) Update(ctx context.Context, tx *sql.Tx, usr models.User) error {

	for i, user := range u.Users {
		if user.Id == usr.Id {
			u.Users[i] = usr
		}
	}

	return nil
}

func (u *InMemoryUsers) Delete(ctx context.Context, _ *sql.Tx, id int) error {
	for i, user := range u.Users {
		if user.Id == id {
			u.Users = append(u.Users[:i], u.Users[i+1:]...)
			break
		}
	}

	return nil
}
