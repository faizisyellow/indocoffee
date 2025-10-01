package service

import (
	"context"
	"database/sql"
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

const CHECK_CONSTRAINT_CART_QUANTITY_CODE = "Check constraint 'cart_items_chk_1'"

var (
	ErrConflictItemCart     = errors.New("carts: item already exist")
	ErrInternalCart         = errors.New("carts: encounter internal error")
	ErrCartNotFound         = errors.New("carts: cart not found")
	ErrCartOverflowQuantity = errors.New("carts: item quantity max is 50")
	ErrCartMinQuantity      = errors.New("carts: item quantity min is 1")
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

func (c *CartsService) FindById(ctx context.Context, id int) (models.Cart, error) {

	cart, err := c.CartsStore.GetById(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.Cart{}, errorService.New(ErrCartNotFound, err)
		default:
			return models.Cart{}, errorService.New(ErrInternalCart, err)
		}
	}

	return cart, nil
}

func (c *CartsService) IncrementItem(ctx context.Context, cartId int) error {
	cart, err := c.CartsStore.GetById(ctx, cartId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return errorService.New(ErrCartNotFound, err)
		default:
			return errorService.New(ErrInternalCart, err)
		}
	}

	if err := c.CartsStore.IncrementQuantity(ctx, cart.Id); err != nil {
		if strings.Contains(err.Error(), CHECK_CONSTRAINT_CART_QUANTITY_CODE) {
			return errorService.New(ErrCartOverflowQuantity, err)
		}
		return errorService.New(ErrInternalCart, err)
	}

	return nil
}

func (c *CartsService) DecrementItem(ctx context.Context, cartId int) error {
	cart, err := c.CartsStore.GetById(ctx, cartId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return errorService.New(ErrCartNotFound, err)
		default:
			return errorService.New(ErrInternalCart, err)
		}
	}

	if err := c.CartsStore.DecrementQuantity(ctx, cart.Id); err != nil {
		if strings.Contains(err.Error(), CHECK_CONSTRAINT_CART_QUANTITY_CODE) {
			return errorService.New(ErrCartOverflowQuantity, err)
		}
		return errorService.New(ErrInternalCart, err)
	}

	return nil

}
