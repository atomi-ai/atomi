package models

type TaxRate struct {
	BaseModel
	Csv                   string  `gorm:"column:csv" json:"csv"`
	State                 string  `gorm:"column:tax_state" json:"state"`
	ZipCode               string  `gorm:"column:zip_code" json:"zip_code"`
	EstimatedCombinedRate float64 `json:"estimated_combined_rate"`
}
