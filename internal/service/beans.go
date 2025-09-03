package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/faizisyellow/indocoffee/internal/repository"
	serviceParser "github.com/faizisyellow/indocoffee/internal/service/parser"
)

type BeansServices struct {
	Repository repository.Repository
}

type RequestCreateBean struct {
	Name string `json:"name" validate:"required,min=3,max=18"`
}

func (rcb RequestCreateBean) Serialize() RequestCreateBean {

	rcb.Name = strings.ToLower(rcb.Name)
	return rcb
}

type RequestUpdateBean struct {
	Name string `json:"name" validate:"required,min=3,max=18"`
}

var (
	ErrConflictBean = errors.New("beans: bean already exist")
	ErrInternalBean = errors.New("beans: encountered an internal error")
	ErrNotFoundBean = errors.New("beans: no such as bean")
)

func (Beans *BeansServices) Create(ctx context.Context, req RequestCreateBean) (string, error) {

	var newBean repository.BeansModel
	newBean.Name = req.Name

	err := Beans.Repository.Beans.Insert(ctx, newBean)
	if err != nil {

		// TODO: should success create new bean if the existing bean is deleted
		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return "", ErrConflictBean
		}
		return "", ErrInternalBean
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
		return ErrInternalBean
	}

	return nil
}

func (Beans *BeansServices) FindAll(ctx context.Context) ([]ResponseFindAll, error) {

	beans, err := Beans.Repository.Beans.GetAll(ctx)
	if err != nil {
		return nil, ErrInternalBean
	}

	response := make([]ResponseFindAll, 0)

	for _, bean := range beans {
		buffRes := new(ResponseFindAll)

		err := serviceParser.Parse(buffRes, bean)
		if err != nil {
			return nil, ErrInternalBean
		}

		response = append(response, *buffRes)
	}

	return response, nil
}

func (Beans *BeansServices) FindById(ctx context.Context, id int) (repository.BeansModel, error) {

	bean, err := Beans.Repository.Beans.GetById(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return repository.BeansModel{}, ErrNotFoundBean
		default:
			return repository.BeansModel{}, ErrInternalBean
		}
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
		// TODO: should success update bean if the existing bean is deleted
		return ErrInternalBean
	}

	return nil
}

func (Beans *BeansServices) Delete(ctx context.Context, id int) error {

	bean, err := Beans.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := Beans.Repository.Beans.Delete(ctx, bean.Id); err != nil {
		return ErrInternalBean
	}

	return nil
}

func (Beans *BeansServices) Remove(ctx context.Context, id int) error {

	bean, err := Beans.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := Beans.Repository.Beans.Destroy(ctx, bean.Id); err != nil {
		return ErrInternalBean
	}

	return nil
}
