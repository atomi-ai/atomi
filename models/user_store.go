package models

// TODO(lamuguo): 用is_enable的方式好像有点奇怪，感觉不用的时候直接删掉更正确一些。
type UserStore struct {
	BaseModel
	User     *User  `gorm:"foreignKey:UserID" json:"-"`
	UserID   int64  `json:"user_id"`
	Store    *Store `gorm:"foreignKey:StoreID" json:"-"`
	StoreID  int64  `json:"store_id"`
	IsEnable bool   `gorm:"column:is_enable" json:"is_enable"`
}
