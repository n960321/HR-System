package model

import "time"

type Account struct {
	ID        uint64      `gorm:"column:id;primaryKey;autoIncrement;type:bigint unsigned;not null;comment:id"`
	Account   string      `gorm:"column:account;type:varchar(255);not null;comment:account"`
	Type      AccountType `gorm:"column:type;type:tinyint unsigned;not null;comment:'1:admin 2:employee'"`
	Name      string      `gorm:"column:name;type:varchar(255);not null;comment:user name"`
	Password  string      `gorm:"column:password;type:varchar(255);not null;comment:encrypted password"`
	CreatedAt time.Time   `gorm:"column:created_at;type:timestamp;null;default:CURRENT_TIMESTAMP;comment:create time"`
	UpdatedAt time.Time   `gorm:"column:updated_at;type:timestamp;null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:update time"`
}

func (Account) TableName() string {
	return "account"
}

type AccountType uint8

const (
	AccountTypeUnknow AccountType = iota
	AccountTypeAdmin
	AccountTypeEmployee
)
