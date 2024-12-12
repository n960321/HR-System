package model

import "time"

type ClockInRecord struct {
	ID        uint64     `gorm:"column:id;primaryKey;autoIncrement;type:bigint unsigned;not null;comment:id"`
	AccountID uint64     `gorm:"column:account_id;type:bigint unsigned;not null"`
	Type      uint8      `gorm:"column:type;type:tinyint unsigned;not null;comment:'1:clock in 2:clock out'"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp;null;default:CURRENT_TIMESTAMP;comment:create time"`
}

func (ClockInRecord) TableName() string {
	return "clock_in_record"
}
