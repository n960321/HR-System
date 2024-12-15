package test

import (
	"HRSystem/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAccount(db *gorm.DB) {
	accounts := []*model.Account{
		{
			Name:     "test1",
			Password: GetPwd("test1"),
			Account:  "test1",
			Type:     model.AccountTypeEmployee,
		},
		{
			Name:     "test2",
			Password: GetPwd("test2"),
			Account:  "test2",
			Type:     model.AccountTypeEmployee,
		},
	}

	for _, a := range accounts {
		db.Save(a)
	}
}

func GetPwd(pwd string) string {
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hashedPwd)
}
