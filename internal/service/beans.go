package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/faizisyellow/indocoffee/internal/repository"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
)

type BeansServices struct {
	Repository repository.Repository
}

const (
	SUCCESS_CREATE_BEAN_MESSAGE = "success create new bean"
)

var (
	ErrConflictBean = errors.New("beans: bean already exist")
	ErrInternalBean = errors.New("beans: encountered an internal error")
	ErrNotFoundBean = errors.New("beans: no such as bean")
)

func (Beans *BeansServices) Create(ctx context.Context, req dto.CreateBeanRequest) (string, error) {

	newBean := repository.BeansModel{
		Name: req.Name,
	}

	err := Beans.Repository.Beans.Insert(ctx, newBean)
	if err != nil {

		// TODO: should success create new bean if the existing bean is deleted
		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return "", errorService.New(ErrConflictBean, err)
		}
		return "", errorService.New(ErrInternalBean, err)
	}

	return SUCCESS_CREATE_BEAN_MESSAGE, nil
}

func (Beans *BeansServices) FindAll(ctx context.Context) ([]repository.BeansModel, error) {

	beans, err := Beans.Repository.Beans.GetAll(ctx)
	if err != nil {
		return nil, errorService.New(ErrInternalBean, err)
	}

	return beans, nil
}

func (Beans *BeansServices) FindById(ctx context.Context, id int) (repository.BeansModel, error) {

	bean, err := Beans.Repository.Beans.GetById(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return repository.BeansModel{}, errorService.New(ErrNotFoundBean, err)
		default:
			return repository.BeansModel{}, errorService.New(ErrInternalBean, err)
		}
	}

	return bean, nil
}

func (Beans *BeansServices) Update(ctx context.Context, id int, req dto.UpdateBeanRequest) error {

	bean, err := Beans.FindById(ctx, id)
	if err != nil {
		return err
	}
	bean.Name = req.Name

	err = Beans.Repository.Beans.Update(ctx, bean)
	if err != nil {
		// TODO: should success update bean if the existing bean is deleted
		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return errorService.New(ErrConflictBean, err)
		}

		return errorService.New(ErrInternalBean, err)
	}

	return nil
}

func (Beans *BeansServices) Delete(ctx context.Context, id int) error {

	bean, err := Beans.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := Beans.Repository.Beans.Delete(ctx, bean.Id); err != nil {
		return errorService.New(ErrInternalBean, err)
	}

	return nil
}

func (Beans *BeansServices) Remove(ctx context.Context) error {

	err := Beans.Repository.Beans.DestroyMany(ctx)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return errorService.New(ErrNotFoundBean, err)
		default:
			return errorService.New(ErrInternalBean, err)
		}
	}

	return nil
}
