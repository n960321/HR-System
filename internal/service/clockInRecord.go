package service

import (
	"HRSystem/internal/model"
	"HRSystem/pkg/database"
)

type ClockInRecordService struct {
	db *database.Database
}

func NewClockInRecordService(db *database.Database) *ClockInRecordService {
	db.GetGorm().AutoMigrate(&model.ClockInRecord{})
	return &ClockInRecordService{
		db: db,
	}
}
