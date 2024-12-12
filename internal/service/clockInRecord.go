package service

import (
	"HRSystem/internal/model"
	"HRSystem/pkg/database"
	"HRSystem/pkg/errors"
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
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

func (s *ClockInRecordService) CreateClockInRecord(ctx context.Context, accountID uint64, clockInType model.ClockInType) error {
	// 查詢最後一筆打卡記錄
	lastRecord := new(model.ClockInRecord)
	err := s.db.GetGorm().Where("account_id = ?", accountID).
		Order("created_at DESC").
		First(lastRecord).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	// 如果是上班打卡
	if clockInType == model.ClockInTypeClockIn {
		// 如果有記錄且最後一筆是上班卡,則返回錯誤
		if err != gorm.ErrRecordNotFound && lastRecord.Type == model.ClockInTypeClockIn {
			return errors.ErrAlreadyClockInToday
		}
	} else if clockInType == model.ClockInTypeClockOut {
		// 如果沒有記錄或最後一筆不是上班卡,則返回錯誤
		if err == gorm.ErrRecordNotFound || lastRecord.Type != model.ClockInTypeClockIn {
			return errors.ErrPleaseClockInFirst
		}
	}

	// 新增一筆打卡記錄
	record := model.ClockInRecord{
		AccountID: accountID,
		Type:      clockInType,
		CreatedAt: time.Now(),
	}
	if err := s.db.GetGorm().Create(&record).Error; err != nil {
		return err
	}
	return nil
}

func (s *ClockInRecordService) ListClockInRecord(ctx context.Context, accountID uint64, start, end time.Time) ([]*model.ClockInRecord, error) {
	records := make([]*model.ClockInRecord, 0)
	log.Debug().Any("accountID", accountID).Any("start", start).Any("end", end).Send()
	err := s.db.GetGorm().Debug().Where("account_id = ? AND created_at BETWEEN ? AND ?", accountID, start, end).
		Order("created_at DESC").
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}
