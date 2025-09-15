package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/mock"
	"github.com/faizisyellow/indocoffee/internal/repository"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
)

func TestCreateBeans(t *testing.T) {
	beansSrv := setupBeansServiceTest()

	t.Run("create a new bean", func(t *testing.T) {
		req := dto.CreateBeanRequest{
			Name: "arabica",
		}
		want := SUCCESS_CREATE_BEAN_MESSAGE
		result, err := beansSrv.Create(context.Background(), req)
		assertNoError(t, err)

		assertBean(t, want, result)
	})

	t.Run("fails create a new bean because it's already exist", func(t *testing.T) {
		req := dto.CreateBeanRequest{
			Name: "robusta",
		}
		_, err := beansSrv.Create(context.Background(), req)
		assertError(t, err, ErrConflictBean.Error())
	})
}

func TestGetBeans(t *testing.T) {
	beansSrv := setupBeansServiceTest()

	t.Run("get bean by id", func(t *testing.T) {
		want := repository.BeansModel{
			Id:   1,
			Name: "arabica",
		}
		reslt, err := beansSrv.FindById(context.Background(), 1)
		assertNoError(t, err)
		assertBean(t, want, reslt)
	})

	t.Run("fails because bean not found by the given id", func(t *testing.T) {
		_, err := beansSrv.FindById(context.Background(), 2)
		assertError(t, err, ErrNotFoundBean.Error())
	})

	t.Run("get all beans", func(t *testing.T) {
		want := []repository.BeansModel{
			{
				Id:   1,
				Name: "arabica",
			},
		}
		reslt, err := beansSrv.FindAll(context.Background())
		assertNoError(t, err)
		assertBean(t, want, reslt)
	})
}

func TestUpdateBeans(t *testing.T) {
	beansSrv := setupBeansServiceTest()

	t.Run("update beans's name by given an id", func(t *testing.T) {
		req := dto.UpdateBeanRequest{Name: "luwak coffee"}
		err := beansSrv.Update(context.Background(), 1, req)
		assertNoError(t, err)
	})

	t.Run("fails update bean because the given bean's id not exist", func(t *testing.T) {
		err := beansSrv.Update(context.Background(), 18, dto.UpdateBeanRequest{})
		assertError(t, err, ErrNotFoundBean.Error())
	})

	t.Run("fails update bean because name already exist", func(t *testing.T) {
		req := dto.UpdateBeanRequest{Name: "robusta"}
		err := beansSrv.Update(context.Background(), 1, req)
		assertError(t, err, ErrConflictBean.Error())
	})
}

func TestDeleteBeans(t *testing.T) {
	beansSrv := setupBeansServiceTest()

	t.Run("delete bean all beans ", func(t *testing.T) {
		err := beansSrv.Remove(context.Background())
		assertNoError(t, err)
	})

	t.Run("delete bean by given an id", func(t *testing.T) {
		err := beansSrv.Delete(context.Background(), 1)
		assertNoError(t, err)
	})

	t.Run("fails delete bean because because the given bean's id not exist", func(t *testing.T) {
		err := beansSrv.Delete(context.Background(), 18)
		assertError(t, err, ErrNotFoundBean.Error())
	})
}

func setupBeansServiceTest() BeansServices {

	return BeansServices{
		Repository: repository.Repository{
			Beans: &mock.BeansRepositoryMock{},
		},
	}

}

func assertBean(t testing.TB, want, got any) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expecting to be match with : %v but got :%v", want, got)
	}
}
