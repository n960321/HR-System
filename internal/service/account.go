package service

import (
	"HRSystem/internal/model"
	"HRSystem/pkg/database"
	"HRSystem/pkg/errors"
	"HRSystem/pkg/helper"
	"HRSystem/pkg/jwthelper"
	"context"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func (s *AccountService) Login(ctx context.Context, account, password string) (string, error) {
	var existingAccount model.Account
	if err := s.db.GetGorm().Where("account = ?", account).First(&existingAccount).Error; err == gorm.ErrRecordNotFound {
		return "", errors.ErrAccountOrPasswordIncorrect
	} else if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingAccount.Password), []byte(password)); err != nil {
		return "", errors.ErrAccountOrPasswordIncorrect
	}

	token, err := jwthelper.GenerateJWTToken(existingAccount.ID, existingAccount.Account, existingAccount.Type)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AccountService) ChangePassword(ctx context.Context, account, oldPassword, newPassword, checkNewPassword string) error {

	var existingAccount model.Account
	if err := s.db.GetGorm().Where("account = ?", account).First(&existingAccount).Error; err != nil {
		return errors.ErrAccountOrPasswordIncorrect
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingAccount.Password), []byte(oldPassword)); err != nil {
		return errors.ErrAccountOrPasswordIncorrect
	}

	if newPassword != checkNewPassword {
		return errors.ErrAccountOrPasswordIncorrect
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.db.GetGorm().Model(&model.Account{}).Where("account = ?", account).Update("password", string(hashedNewPassword)).Error; err != nil {
		return err
	}

	return nil

}

func (s *AccountService) CreateAccount(ctx context.Context, name, account string) (string, error) {
	var existingAccount model.Account
	if err := s.db.GetGorm().Where("account = ?", account).First(&existingAccount).Error; err == nil {
		return "", errors.ErrAccountDuplicate
	}

	// 隨機產生10碼的pwd
	pwd := helper.GenerateRandomString(10)

	// 對pwd進行 bcrypt 加密
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// 存進DB
	newAccount := model.Account{
		Account:  account,
		Name:     name,
		Type:     model.AccountTypeEmployee,
		Password: string(hashedPwd),
	}

	if err := s.db.GetGorm().Create(&newAccount).Error; err != nil {
		return "", err
	}

	return pwd, nil
}
