package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository/beans"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
)

type BeansServices struct {
	BeansStore beans.Beans
}

const (
	SUCCESS_CREATE_BEAN_MESSAGE = "success create new bean"
)

var (
	ErrConflictBean = errors.New("beans: bean already exist")
	ErrInternalBean = errors.New("beans: encountered an internal error")
	ErrNotFoundBean = errors.New("beans: no such as bean")
)

func (b *BeansServices) Create(ctx context.Context, req dto.CreateBeanRequest) (string, error) {

	newBean := models.BeansModel{
		Name: req.Name,
	}

	err := b.BeansStore.Insert(ctx, newBean)
	if err != nil {

		// TODO: should success create new bean if the existing bean is deleted
		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return "", errorService.New(ErrConflictBean, err)
		}
		return "", errorService.New(ErrInternalBean, err)
	}

	return SUCCESS_CREATE_BEAN_MESSAGE, nil
}

func (b *BeansServices) FindAll(ctx context.Context) ([]models.BeansModel, error) {

	beans, err := b.BeansStore.GetAll(ctx)
	if err != nil {
		return nil, errorService.New(ErrInternalBean, err)
	}

	return beans, nil
}

func (b *BeansServices) FindById(ctx context.Context, id int) (models.BeansModel, error) {

	bean, err := b.BeansStore.GetById(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.BeansModel{}, errorService.New(ErrNotFoundBean, err)
		default:
			return models.BeansModel{}, errorService.New(ErrInternalBean, err)
		}
	}

	return bean, nil
}

func (b *BeansServices) Update(ctx context.Context, id int, req dto.UpdateBeanRequest) error {

	bean, err := b.FindById(ctx, id)
	if err != nil {
		return err
	}
	bean.Name = req.Name

	err = b.BeansStore.Update(ctx, bean)
	if err != nil {
		// TODO: should success update bean if the existing bean is deleted
		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return errorService.New(ErrConflictBean, err)
		}

		return errorService.New(ErrInternalBean, err)
	}

	return nil
}

func (b *BeansServices) Delete(ctx context.Context, id int) error {

	bean, err := b.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := b.BeansStore.Delete(ctx, bean.Id); err != nil {
		return errorService.New(ErrInternalBean, err)
	}

	return nil
}

func (b *BeansServices) Remove(ctx context.Context) error {

	err := b.BeansStore.DestroyMany(ctx)
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
