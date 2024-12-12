package errors

import "errors"

var (
	// 已經打過上班卡，請先打下班卡
	ErrAlreadyClockInToday = errors.New("already clocked in today, please clock out first")
	// 請先打下班卡，在打下班卡
	ErrPleaseClockInFirst = errors.New("please clock in before clocking out")
)
