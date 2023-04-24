package models

// Config represents the config entity
type Config struct {
	Key   string `gorm:"primaryKey;column:config_key" json:"key"`
	Value string `gorm:"column:config_value" json:"value"`
}
