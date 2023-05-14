package models

// TODO(lamuguo): Merge magager_stores and user_stores as one table,
// because we just need to have a user x store relationship table.
//
// 否则会有的一个问题：譬如 /api/store/{store_id} 的时候，你会很困惑要query
// 哪张表。如果分成 /api/store/{store_id} 和 /api/mgr/store/{store_id}
// 的话，一开始你给mgr / store建联系的时候，你就需要往两张表里面同时插入关系
// 要不总不能你能管这个店了，但是却不能同时看这个店吧。
//
// 所以，比较好的办法是只用一张表，但是有个field标注他们是不是is_managed的关系。
type ManagerStores struct {
	BaseModel
	UserID  int64 `gorm:"uniqueIndex:idx_manager_stores"`
	StoreID int64 `gorm:"uniqueIndex:idx_manager_stores"`
}

func (ManagerStores) TableName() string {
	return "manager_stores"
}
