package service

import (
	"context"

	"github.com/faizisyellow/indocoffee/internal/repository"
)

type BeansServices struct {
	Repository repository.Repository
}

type RequestCreateBean struct {
	Name string `json:"name" validate:"required,min=3,max=18"`
}

type RequestUpdateBean struct {
	Name string `json:"name" validate:"required,min=3,max=18"`
}

func (Beans *BeansServices) Create(ctx context.Context, req RequestCreateBean) (string, error) {

	var newBean repository.BeansModel
	newBean.Name = req.Name

	err := Beans.Repository.Beans.Insert(ctx, newBean)
	if err != nil {
		return "", err
	}

	return "create new bean successfully", nil
}

func (Beans *BeansServices) FindAll(ctx context.Context) ([]repository.BeansModel, error) {

	beans, err := Beans.Repository.Beans.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return beans, nil
}

func (Beans *BeansServices) FindById(ctx context.Context, id int) (repository.BeansModel, error) {

	bean, err := Beans.Repository.Beans.GetById(ctx, id)
	if err != nil {
		return repository.BeansModel{}, err
	}

	return bean, nil
}

func (Beans *BeansServices) Update(ctx context.Context, id int, nw repository.BeansModel) error {

	bean, err := Beans.FindById(ctx, id)
	if err != nil {
		return err
	}

	bean.Name = nw.Name

	err = Beans.Repository.Beans.Update(ctx, bean)
	if err != nil {
		return err
	}

	return nil
}

func (Beans *BeansServices) Delete(ctx context.Context, id int) error {

	bean, err := Beans.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := Beans.Repository.Beans.Delete(ctx, bean.Id); err != nil {
		return err
	}

	return nil
}

func (Beans *BeansServices) Remove(ctx context.Context, id int) error {

	bean, err := Beans.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := Beans.Repository.Beans.Destroy(ctx, bean.Id); err != nil {
		return err
	}

	return nil
}
