package models

import (
	"encoding/json"
)

type PaymentIntentRequest struct {
	Amount            int64         `json:"amount"`
	Currency          string        `json:"currency"`
	PaymentMethodID   string        `json:"payment_method_id"`
	ShippingAddressID int64         `json:"shipping_address_id"`
	OrderID           int64         `json:"order_id"`
	DeliveryData      *DeliveryData `json:"delivery_data"`
}

// UnmarshalJSONPaymentIntentRequest 将JSON字符串转换为PaymentIntentRequest结构体
func UnmarshalJSONPaymentIntentRequest(data []byte) (*PaymentIntentRequest, error) {
	var req PaymentIntentRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}
