package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Balance uint `json:"balance" gorm:"default:0"`
}

type Reserve struct {
	gorm.Model
	UserID         uint `gorm:"column:user_id"`
	ServiceID      uint
	OrderID        uint
	Price          uint
	ReservedAmount uint
}

// ReportEntry - используется для создания отчетов.
// Должно быть заполнено при каждом признании выручки, включая идентификатор пользователя, идентификатор заказа и фактическую выручку.
type ReportEntry struct {
	gorm.Model
	UserID  uint
	OrderID uint
	Revenue uint
}
