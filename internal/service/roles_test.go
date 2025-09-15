package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/mock"
	"github.com/faizisyellow/indocoffee/internal/repository"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	"github.com/faizisyellow/indocoffee/internal/utils"
)

func TestCreateRoles(t *testing.T) {
	rolesSrv := setupRolesServiceTest()
	t.Run("create new role", func(t *testing.T) {
		req := dto.CreateRoleRequest{
			Name:  "customer",
			Level: 1,
		}
		want := SUCCESS_CREATE_ROLES_MESSAGE
		got, err := rolesSrv.Create(context.Background(), req)
		assertNoError(t, err)
		if want != got {
			t.Errorf("expecting to be match: %v but got: %v", want, got)
		}
	})

	t.Run("fails create a new role because role already exist", func(t *testing.T) {
		req := dto.CreateRoleRequest{
			Name:  "admin",
			Level: 1,
		}
		_, err := rolesSrv.Create(context.Background(), req)
		assertError(t, err, ErrConflictRole.Error())
	})
}

func TestGetRoles(t *testing.T) {
	rolesSrv := setupRolesServiceTest()

	t.Run("get roles by id", func(t *testing.T) {
		want := repository.RolesModel{
			Id:    1,
			Name:  "admin",
			Level: 3,
		}
		got, err := rolesSrv.FindById(context.Background(), 1)
		assertNoError(t, err)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expecting to be match: %v but got: %v", want, got)
		}
	})

	t.Run("fails get roles because role not found", func(t *testing.T) {
		_, err := rolesSrv.FindById(context.Background(), 2)
		assertError(t, err, ErrNotFoundRole.Error())
	})

	t.Run("get all roles", func(t *testing.T) {
		want := []repository.RolesModel{
			{
				Id:    1,
				Name:  "admin",
				Level: 3,
			},
		}
		got, err := rolesSrv.FindAll(context.Background())
		assertNoError(t, err)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expecting to be match: %v but got: %v", want, got)
		}
	})
}

func TestUpdateRole(t *testing.T) {
	rolesSrv := setupRolesServiceTest()

	t.Run("update only the name of role", func(t *testing.T) {
		req := dto.UpdateRoleRequest{
			Name: utils.StringToPoint("super admin"),
		}
		err := rolesSrv.Update(context.Background(), 1, req)
		assertNoError(t, err)
	})

	t.Run("update only the level of role", func(t *testing.T) {
		req := dto.UpdateRoleRequest{
			Level: utils.IntToPoint(4),
		}
		err := rolesSrv.Update(context.Background(), 1, req)
		assertNoError(t, err)
	})

	t.Run("fails update role because role not found", func(t *testing.T) {
		err := rolesSrv.Update(context.Background(), 10, dto.UpdateRoleRequest{})
		assertError(t, err, ErrNotFoundRole.Error())
	})

	t.Run("fails update role because the name already exist", func(t *testing.T) {
		err := rolesSrv.Update(context.Background(), 1, dto.UpdateRoleRequest{
			Name: utils.StringToPoint("manager"),
		})
		assertError(t, err, ErrConflictRole.Error())
	})

}

func TestDeleteRoles(t *testing.T) {
	rolesSrv := setupRolesServiceTest()

	t.Run("delete all roles", func(t *testing.T) {
		err := rolesSrv.Remove(context.Background())
		assertNoError(t, err)
	})

	t.Run("delete roles by id", func(t *testing.T) {
		err := rolesSrv.Delete(context.Background(), 1)
		assertNoError(t, err)
	})

	t.Run("fails delete roles because not such as role by given id", func(t *testing.T) {
		err := rolesSrv.Delete(context.Background(), 10)
		assertError(t, err, ErrNotFoundRole.Error())
	})

}

func setupRolesServiceTest() RolesServices {
	return RolesServices{
		Repository: repository.Repository{
			Roles: &mock.RolesRepositoryMock{},
		},
	}
}
