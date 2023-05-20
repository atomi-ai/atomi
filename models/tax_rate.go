package models

type TaxRate struct {
	BaseModel
	Csv                   string  `gorm:"column:csv" json:"csv"`
	State                 string  `gorm:"column:tax_state" json:"state"`
	ZipCode               string  `gorm:"column:zip_code" json:"zip_code"`
	TaxRegionName         string  `gorm:"column:tax_region_name" json:"tax_region_name"`
	EstimatedCombinedRate float64 `json:"estimated_combined_rate"`
	StateRate             float64 `json:"state_rate"`
	EstimatedCountyRate   float64 `json:"estimated_county_rate"`
	EstimatedCityRate     float64 `json:"estimated_city_rate"`
	EstimatedSpecialRate  float64 `json:"estimated_special_rate"`
	RiskLevel             int     `gorm:"column:risk_level" json:"risk_level"`
}
