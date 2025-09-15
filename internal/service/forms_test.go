package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/mock"
	"github.com/faizisyellow/indocoffee/internal/repository"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
)

func TestCreateForms(t *testing.T) {
	formsSrv := setupFormsServiceTest()

	t.Run("create a new form", func(t *testing.T) {
		req := dto.CreateFormRequest{
			Name: "whole coffee beans",
		}
		want := SUCCESS_CREATE_FORMS_MESSAGE
		resl, err := formsSrv.Create(context.Background(), req)
		assertNoError(t, err)
		if want != resl {
			t.Errorf("expecting to be match with : %v but got : %v", want, resl)
		}
	})

	t.Run("fails cretae a new form because form already exist", func(t *testing.T) {
		req := dto.CreateFormRequest{
			Name: "grounded",
		}

		_, err := formsSrv.Create(context.Background(), req)
		assertError(t, err, ErrConflictForm.Error())
	})
}

func TestGetForms(t *testing.T) {
	formsSrv := setupFormsServiceTest()

	t.Run("get forms by id", func(t *testing.T) {
		want := repository.FormsModel{
			Id:   1,
			Name: "grounded",
		}
		got, err := formsSrv.FindById(context.Background(), want.Id)
		assertNoError(t, err)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected to be match: %v but got: %v", want, got)
		}
	})

	t.Run("fauls get forms by id because form not found", func(t *testing.T) {

		_, err := formsSrv.FindById(context.Background(), 2)
		assertError(t, err, ErrNotFoundForm.Error())

	})

	t.Run("get all forms", func(t *testing.T) {
		want := []repository.FormsModel{
			{Id: 1, Name: "grounded"},
		}
		got, err := formsSrv.FindAll(context.Background())
		assertNoError(t, err)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected to be match: %v but got: %v", want, got)
		}
	})
}

func TestUpdateForms(t *testing.T) {
	formsSrv := setupFormsServiceTest()

	t.Run("update form name by given an id", func(t *testing.T) {
		req := dto.UpdateFormRequest{
			Name: "something form",
		}
		err := formsSrv.Update(context.Background(), 1, req)
		assertNoError(t, err)
	})

	t.Run("fails update form because form not found", func(t *testing.T) {
		err := formsSrv.Update(context.Background(), 34, dto.UpdateFormRequest{})
		assertError(t, err, ErrNotFoundForm.Error())
	})

	t.Run("fails update form because the name of the form already exist", func(t *testing.T) {
		req := dto.UpdateFormRequest{
			Name: "grounded",
		}
		err := formsSrv.Update(context.Background(), 1, req)
		assertError(t, err, ErrConflictForm.Error())
	})
}

func TestDeleteForms(t *testing.T) {
	formsSrv := setupFormsServiceTest()

	t.Run("delete all forms", func(t *testing.T) {
		err := formsSrv.Remove(context.Background())
		assertNoError(t, err)
	})

	t.Run("delete a form by an id", func(t *testing.T) {
		err := formsSrv.Delete(context.Background(), 1)
		assertNoError(t, err)
	})

	t.Run("fails delete a form because form not found", func(t *testing.T) {
		err := formsSrv.Delete(context.Background(), 20)
		assertError(t, err, ErrNotFoundForm.Error())
	})
}

func setupFormsServiceTest() FormsServices {
	return FormsServices{
		Repository: repository.Repository{
			Forms: &mock.FormsRepositoryMock{},
		},
	}
}
