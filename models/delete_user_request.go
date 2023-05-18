package models

type DeleteUserRequest struct {
	BaseModel
	UserID int64 `gorm:"unique" json:"user_id"`
}
