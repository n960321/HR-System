package test

import (
	"HRSystem/internal/model"
	"HRSystem/pkg/errors"
	"context"
	"testing"
)

func TestAccountService_Login(t *testing.T) {
	type args struct {
		ctx      context.Context
		account  string
		password string
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Admin Login",
			args: args{
				ctx:      context.Background(),
				account:  "admin",
				password: "123456",
			},
			err: nil,
		},
		{
			name: "General Account Login",
			args: args{
				ctx:      context.Background(),
				account:  "test1",
				password: "test1",
			},
			err: nil,
		},
		{
			name: "Login with Nonexistent Account",
			args: args{
				ctx:      context.Background(),
				account:  "nonexisttest",
				password: "nonexisttest",
			},
			err: errors.ErrAccountOrPasswordIncorrect,
		},
		{
			name: "Login with Incorrect Password",
			args: args{
				ctx:      context.Background(),
				account:  "test1",
				password: "wrongpassword",
			},
			err: errors.ErrAccountOrPasswordIncorrect,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := accountSvc.Login(tt.args.ctx, tt.args.account, tt.args.password)
			if err != tt.err {
				t.Errorf("AccountService.Login() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestAccountService_ChangePassword(t *testing.T) {
	type args struct {
		ctx              context.Context
		account          string
		oldPassword      string
		newPassword      string
		checkNewPassword string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Correctly Change Password",
			args: args{
				ctx:              context.Background(),
				account:          "test1",
				oldPassword:      "test1",
				newPassword:      "test1new",
				checkNewPassword: "test1new",
			},
			wantErr: nil,
		},
		{
			name: "Change Password with Nonexistent Account",
			args: args{
				ctx:              context.Background(),
				account:          "notexist",
				oldPassword:      "notexist",
				newPassword:      "notexist",
				checkNewPassword: "notexist",
			},
			wantErr: errors.ErrAccountOrPasswordIncorrect,
		},
		{
			name: "Old Password Input Error",
			args: args{
				ctx:              context.Background(),
				account:          "test1",
				oldPassword:      "wrongoldpassword",
				newPassword:      "newpassword",
				checkNewPassword: "newpassword",
			},
			wantErr: errors.ErrAccountOrPasswordIncorrect,
		},
		{
			name: "New Password and Confirm New Password Do Not Match",
			args: args{
				ctx:              context.Background(),
				account:          "test2",
				oldPassword:      "test2",
				newPassword:      "newpassword",
				checkNewPassword: "differentpassword",
			},
			wantErr: errors.ErrAccountOrPasswordIncorrect,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := accountSvc.ChangePassword(tt.args.ctx, tt.args.account, tt.args.oldPassword, tt.args.newPassword, tt.args.checkNewPassword); err != tt.wantErr {
				t.Errorf("AccountService.ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccountService_CreateAccount(t *testing.T) {
	type args struct {
		ctx     context.Context
		creator *model.Account
		name    string
		account string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Successfully Create Account",
			args: args{
				ctx:     context.Background(),
				creator: &model.Account{Type: model.AccountTypeAdmin},
				name:    "create_test3",
				account: "create_test3",
			},
			wantErr: nil,
		},
		{
			name: "Creator Privilege Insufficient",
			args: args{
				ctx:     context.Background(),
				creator: &model.Account{Type: model.AccountTypeEmployee},
				name:    "test4",
				account: "test4",
			},
			wantErr: errors.ErrInsufficientPrivilege,
		},
		{
			name: "Create Duplicate Account",
			args: args{
				ctx:     context.Background(),
				creator: &model.Account{Type: model.AccountTypeAdmin},
				name:    "test2",
				account: "test2",
			},
			wantErr: errors.ErrAccountDuplicate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := accountSvc.CreateAccount(tt.args.ctx, tt.args.creator, tt.args.name, tt.args.account)
			if err != tt.wantErr {
				t.Errorf("AccountService.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
