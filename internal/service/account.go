package service

import (
	"HRSystem/internal/model"
	"HRSystem/pkg/database"
	"context"
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

func (s *AccountService) Login(ctx context.Context) error {
	panic("not implement")
}

func (s *AccountService) ChangePassword(ctx context.Context) {
	panic("not implement")
}

func (s *AccountService) CreateAccount(ctx context.Context) {
	panic("not implement")
}
