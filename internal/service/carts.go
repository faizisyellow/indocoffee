package service

import (
	"context"
	"errors"
	"strings"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository/carts"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
)

type CartsService struct {
	CartsStore      carts.Carts
	ProductsService ProductsServiceInterface
}

var (
	ErrConflictItemCart = errors.New("carts: item already exist")
	ErrInternalCart     = errors.New("carts: encounter internal error")
)

func (c *CartsService) Create(ctx context.Context, req dto.CreateCartRequest, usrId int) error {
	prd, err := c.ProductsService.FindById(ctx, req.ProductId)
	if err != nil {
		return err
	}

	if err := c.CartsStore.Insert(ctx, models.Cart{
		ProductId: prd.Id,
		UserId:    usrId,
	}); err != nil {
		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return errorService.New(ErrConflictItemCart, err)
		}
		return errorService.New(ErrInternalCart, err)
	}

	return nil
}
