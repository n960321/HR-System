package errors

import "errors"

var (
	// 已經打過上班卡，請先打下班卡
	ErrAlreadyClockInToday = errors.New("already clocked in today, please clock out first")
	// 請先打下班卡，在打下班卡
	ErrPleaseClockInFirst = errors.New("please clock in before clocking out")
	// 帳號已重複
	ErrAccountDuplicate = errors.New("account already exists")
	// 帳號或密碼錯誤
	ErrAccountOrPasswordIncorrect = errors.New("account or password is incorrect")
	// 不合法的 token
	ErrInvalidToken = errors.New("invalid token")
)
