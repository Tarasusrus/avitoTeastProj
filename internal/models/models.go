package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Balance uint `json:"balance" gorm:"default:0"`
}

type Service struct {
	gorm.Model
	Price uint
	Name  string
}

// Order - представляет заказ, связывается с Service для определения стоимости и с User для указания, кто сделал заказ.
// Поле Status может быть использовано для отслеживания состояния заказа, например "reserved" при резервировании баланса, и "confirmed" при списании средств и признании выручки.
type Order struct {
	gorm.Model
	Service uint
	User    uint
	Cost    uint
	Status  string // "reserved", "confirmed",...
}

type Reserve struct {
	gorm.Model
	UserID         uint
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
