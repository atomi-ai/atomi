package models

import (
	"time"
)

type BaseModel struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO(lamuguo): 这个的delete不会真的delete，只是mark了deletedAt，所以会有一些问题。暂时disable，以后再enable。
	//DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
