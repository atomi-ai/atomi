package models

type Order struct {
	BaseModel
	UserID          int64       `gorm:"column:user_id" json:"user_id"`
	StoreID         int64       `gorm:"column:store_id" json:"store_id"`
	PaymentIntentID *string     `gorm:"column:payment_intent_id;unique" json:"payment_intent_id"`
	DeliveryID      *string     `gorm:"column:delivery_id;unique" json:"delivery_id"`
	OrderItems      []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
	DisplayStatus   OrderStatus `gorm:"column:status" json:"display_status"`
}

type OrderStatus string

const (
	OrderStatusWaitingForPayment = OrderStatus("WAITING_FOR_PAYMENT")
	OrderStatusPaid              = OrderStatus("PAID")
	OrderStatusInProduction      = OrderStatus("IN_PRODUCTION")
	OrderStatusReadyForPickup    = OrderStatus("READY_FOR_PICKUP")
	OrderStatusInDelivery        = OrderStatus("IN_DELIVERY")
	OrderStatusCompleted         = OrderStatus("COMPLETED")
	OrderStatusRefunded          = OrderStatus("REFUNDED")
)

var DeliveryToOrderStatus = map[DeliveryStatus]OrderStatus{
	DeliveryStatusPending:        OrderStatusReadyForPickup,
	DeliveryStatusPickup:         OrderStatusReadyForPickup,
	DeliveryStatusPickupComplete: OrderStatusInDelivery,
	DeliveryStatusDropoff:        OrderStatusInDelivery,
	DeliveryStatusDelivered:      OrderStatusCompleted,
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	BaseModel
	OrderID   int64    `gorm:"column:order_id" json:"order_id"`
	Product   *Product `gorm:"foreignKey:ProductID" json:"product"`
	ProductID int64    `gorm:"column:product_id" json:"product_id"`
	Quantity  int64    `gorm:"column:quantity" json:"quantity"`
}
