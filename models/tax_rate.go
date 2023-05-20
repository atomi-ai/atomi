package models

type TaxRate struct {
	BaseModel
	Csv                   string  `gorm:"column:csv;type:varchar(255)" json:"csv"`
	State                 string  `gorm:"column:State;type:char(2)" json:"State"`
	ZipCode               string  `gorm:"column:ZipCode;type:char(5)" json:"ZipCode"`
	TaxRegionName         string  `gorm:"column:TaxRegionName;type:varchar(255)" json:"TaxRegionName"`
	EstimatedCombinedRate float64 `gorm:"column:EstimatedCombinedRate;type:decimal(8,6)" json:"EstimatedCombinedRate"`
	StateRate             float64 `gorm:"column:StateRate;type:decimal(8,6)" json:"StateRate"`
	EstimatedCountyRate   float64 `gorm:"column:EstimatedCountyRate;type:decimal(8,6)" json:"EstimatedCountyRate"`
	EstimatedCityRate     float64 `gorm:"column:EstimatedCityRate;type:decimal(8,6)" json:"EstimatedCityRate"`
	EstimatedSpecialRate  float64 `gorm:"column:EstimatedSpecialRate;type:decimal(8,6)" json:"EstimatedSpecialRate"`
	RiskLevel             int     `gorm:"column:RiskLevel;type:int" json:"RiskLevel"`
}
