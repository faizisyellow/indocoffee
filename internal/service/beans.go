package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

var (
	ErrConflictBean = errors.New("beans: bean already exist")
)

func (Beans *BeansServices) Create(ctx context.Context, req RequestCreateBean) (string, error) {

	var newBean repository.BeansModel
	newBean.Name = req.Name

	err := Beans.Repository.Beans.Insert(ctx, newBean)
	if err != nil {
		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return "", ErrConflictBean
		}
		return "", err
	}

	return "create new bean successfully", nil
}

type ResponseFindAll struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (rf *ResponseFindAll) ParseDTO(data any) error {

	switch v := data.(type) {
	case repository.BeansModel:
		rf.Id = v.Id
		rf.Name = v.Name
	default:
		return fmt.Errorf("something went wrong when response to client")
	}

	return nil
}

func (Beans *BeansServices) FindAll(ctx context.Context) ([]ResponseFindAll, error) {

	beans, err := Beans.Repository.Beans.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]ResponseFindAll, 0)

	for _, bean := range beans {
		buffRes := new(ResponseFindAll)

		err := buffRes.ParseDTO(bean)
		if err != nil {
			return nil, err
		}

		response = append(response, *buffRes)
	}

	return response, nil
}

func (Beans *BeansServices) FindById(ctx context.Context, id int) (repository.BeansModel, error) {

	bean, err := Beans.Repository.Beans.GetById(ctx, id)
	if err != nil {
		return repository.BeansModel{}, err
	}

	return bean, nil
}

func updateBeanPayload(req RequestUpdateBean, ltsBean repository.BeansModel) repository.BeansModel {

	ltsBean.Name = req.Name

	return ltsBean
}

func (Beans *BeansServices) Update(ctx context.Context, id int, req RequestUpdateBean) error {

	bean, err := Beans.FindById(ctx, id)
	if err != nil {
		return err
	}

	err = Beans.Repository.Beans.Update(ctx, updateBeanPayload(req, bean))
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
