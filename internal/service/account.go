package service

import (
	"HRSystem/internal/model"
	"HRSystem/pkg/database"
)

type AccountService struct {
	db *database.Database
}

func NewAccountService(db *database.Database) *AccountService {
	db.GetGorm().AutoMigrate(&model.Account{})
	return &AccountService{
		db: db,
	}
}
