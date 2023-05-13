package services

import (
	"errors"

	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
	"gorm.io/gorm"
)

type OrderService interface {
	GetUserOrders(userID int64) ([]models.Order, error)
	AddOrderForUser(user *models.User, order *models.Order) (*models.Order, error)
	FindOrderByID(orderID int64) (*models.Order, error)
	UpdatePaymentIntentID(orderID int64, paymentIntentID string) (*models.Order, error)
	UpdateDeliveryID(orderID int64, deliveryID string) (*models.Order, error)
}

type orderService struct {
	OrderRepo     repositories.OrderRepository
	OrderItemRepo repositories.OrderItemRepository
	StripeService StripeService
}

func NewOrderService(orderRepo repositories.OrderRepository, orderItemRepo repositories.OrderItemRepository, stripeService StripeService) OrderService {
	return &orderService{
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		StripeService: stripeService,
	}
}

func (os *orderService) GetUserOrders(userID int64) ([]models.Order, error) {
	orders, err := os.OrderRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	for i := range orders {
		if orders[i].PaymentIntentID == nil || *orders[i].PaymentIntentID == "" {
			orders[i].DisplayStatus = "pending payment"
			continue
		}

		paymentIntent, err := os.StripeService.RetrievePaymentIntent(*orders[i].PaymentIntentID)
		if err != nil {
			return nil, err
		}

		refunded := paymentIntent.LatestCharge.Refunded
		if refunded {
			orders[i].DisplayStatus = "refunded"
		} else {
			orders[i].DisplayStatus = string(paymentIntent.Status)
		}
	}

	return orders, nil
}

func processOrderItems(orderItems []models.OrderItem) {
	for i := range orderItems {
		if orderItems[i].Product != nil {
			orderItems[i].ProductID = orderItems[i].Product.ID
		}
	}
}

func (os *orderService) AddOrderForUser(user *models.User, order *models.Order) (*models.Order, error) {
	if order.UserID == 0 {
		order.UserID = user.ID
	}

	processOrderItems(order.OrderItems)
	err := os.OrderRepo.Save(order)
	if err != nil {
		return nil, err
	}

	for i := range order.OrderItems {
		order.OrderItems[i].OrderID = order.ID
		err := os.OrderItemRepo.Save(&order.OrderItems[i])
		if err != nil {
			return nil, err
		}
	}

	return order, nil
}

func (os *orderService) FindOrderByID(orderID int64) (*models.Order, error) {
	order, err := os.OrderRepo.GetByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return order, nil
}

func (os *orderService) UpdatePaymentIntentID(orderID int64, paymentIntentID string) (*models.Order, error) {
	order, err := os.OrderRepo.GetByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ORDER NOT FOUND")
		}
		return nil, err
	}

	order.PaymentIntentID = &paymentIntentID
	err = os.OrderRepo.Save(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (os *orderService) UpdateDeliveryID(orderID int64, deliveryID string) (*models.Order, error) {
	order, err := os.OrderRepo.GetByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ORDER NOT FOUND")
		}
		return nil, err
	}

	order.DeliveryID = &deliveryID
	err = os.OrderRepo.Save(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}
