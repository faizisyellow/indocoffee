package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/faizisyellow/indocoffee/internal/db"
	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository/carts"
	"github.com/faizisyellow/indocoffee/internal/repository/orders"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"github.com/faizisyellow/indocoffee/internal/utils"
)

type OrdersService struct {
	UsersService    UsersServiceInterface
	ProductsService ProductsServiceInterface
	CartsStore      carts.Carts
	OrderStore      orders.Orders
	Transaction     db.Transactioner
	Uuid            utils.Token
}

const (
	ORDER_ID_PREFIX = "ORD"
)

var (
	ErrOrdersInternal      = errors.New("orders: encounter internal error")
	ErrOrdersItemEmpty     = errors.New("orders: items empty")
	ErrOrdersConflict      = errors.New("orders: already exist")
	ErrOrdersQuantityIssue = errors.New("orders: one of the item is not available")
)

func (o *OrdersService) Create(ctx context.Context, idemKey string, req dto.CreateOrderRequest, usrId int) error {
	idempotencyKeyExist, err := o.OrderStore.GetIdempotencyKey(ctx, idemKey)
	if err != nil && err != sql.ErrNoRows {
		return errorService.New(ErrOrdersInternal, err)
	}

	if idempotencyKeyExist != "" {
		return errorService.New(ErrOrdersConflict, ErrOrdersConflict)
	}

	customer, err := o.UsersService.FindUserById(ctx, usrId)
	if err != nil {
		return err
	}

	var (
		items            []models.OrderItem
		totalPrice       float64
		alternativePhone *string
		cartItems        []models.Cart
	)

	for _, cartId := range req.CartIds {
		cart, err := o.CartsStore.GetById(ctx, cartId)
		if err != nil {
			// should be already handle by the authorization
			log.Printf("error getting product: %v", err.Error())
		}

		cartItems = append(cartItems, models.Cart{
			Id:        cart.Id,
			ProductId: cart.ProductId,
			UserId:    cart.UserId,
			Quantity:  cart.Quantity,
		})
	}

	for _, item := range cartItems {
		product, err := o.ProductsService.FindById(ctx, item.ProductId)
		if err != nil {
			log.Printf("error getting product: %v", err.Error())
			continue
		}

		if item.Quantity >= product.Quantity || product.Quantity <= 0 {
			return ErrOrdersQuantityIssue
		}

		totalPrice += product.Price * float64(item.Quantity)
		items = append(items, models.OrderItem{
			Id:            product.Id,
			Image:         product.Image,
			BeanName:      product.BeansModel.Name,
			FormName:      product.FormsModel.Name,
			Roasted:       product.Roasted,
			Price:         product.Price,
			OrderQuantity: item.Quantity,
		})
	}

	if len(items) == 0 {
		return errorService.New(ErrOrdersItemEmpty, ErrOrdersItemEmpty)
	}

	if req.AlternativePhoneNumber != nil {
		cleantAlt, err := utils.ValidateAndFormatPhoneNumber(*req.AlternativePhoneNumber)
		if err != nil {
			return err
		}
		alternativePhone = &cleantAlt
	}

	id := strings.Builder{}
	id.WriteString(ORDER_ID_PREFIX)
	id.WriteString("-")
	id.WriteString(o.Uuid.Generate())

	req.PhoneNumber, err = utils.ValidateAndFormatPhoneNumber(req.PhoneNumber)
	if err != nil {
		return err
	}

	newOrder := models.Order{
		Id:             id.String(),
		IdempotencyKey: idemKey,
		// user authenticate
		CustomerId:             customer.Id,
		CustomerName:           req.CustomerName,
		CustomerEmail:          req.CustomerEmail,
		PhoneNumber:            req.PhoneNumber,
		Street:                 req.Street,
		City:                   req.City,
		AlternativePhoneNumber: alternativePhone,
		Items:                  items,
		TotalPrice:             totalPrice,
	}

	return o.Transaction.WithTx(ctx, func(tx *sql.Tx) error {
		if err := o.OrderStore.Create(ctx, tx, newOrder); err != nil {
			return errorService.New(ErrOrdersInternal, err)
		}

		for _, cartId := range req.CartIds {
			err := o.CartsStore.UpdateCartStatus(ctx, tx, cartId, carts.Ordered)
			if err != nil {
				return errorService.New(ErrOrdersInternal, err)
			}
		}

		for _, item := range cartItems {
			err := o.ProductsService.DecreaseQuantityProduct(ctx, tx, item.ProductId, item.Quantity)
			if err != nil {
				if strings.Contains(err.Error(), "quantity_non_negative") {
					return errorService.New(ErrCartMinQuantity, err)
				}
				return errorService.New(ErrOrdersInternal, err)
			}
		}

		return nil
	})
}
