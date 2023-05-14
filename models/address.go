package models

type Address struct {
	BaseModel
	Line1      string `json:"line1"`
	Line2      string `json:"line2"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `gorm:"column:postal_code" json:"postal_code"`
}
