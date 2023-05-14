package models

type UserAddress struct {
	BaseModel
	User      *User    `gorm:"foreignKey:UserID" json:"-"`
	UserID    int64    `json:"user_id"`
	Address   *Address `gorm:"foreignKey:AddressID" json:"-"`
	AddressID int64    `json:"address_id"`
}
