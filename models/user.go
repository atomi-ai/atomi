package models

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
	RoleMgr   Role = "MANAGER"
)

type User struct {
	BaseModel
	Email                    string  `gorm:"unique" json:"email"`
	Role                     Role    `json:"role"`
	Phone                    string  `json:"phone"`
	Name                     string  `json:"name"`
	DefaultShippingAddressID int64   `json:"default_shipping_address_id" gorm:"column:default_shipping_address_id"`
	DefaultBillingAddressID  int64   `json:"default_billing_address_id" gorm:"column:default_billing_address_id"`
	StripeCustomerID         string  `json:"stripe_customer_id" gorm:"column:stripe_customer_id"`
	PaymentMethodID          *string `json:"payment_method_id" gorm:"column:payment_method_id"`
}

func (User) TableName() string {
	return "users"
}
