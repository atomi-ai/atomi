package models

// Config represents the config entity
type Config struct {
	BaseModel
	Key   string `gorm:"column:config_key;unique" json:"key"`
	Value string `gorm:"column:config_value" json:"value"`
}
