package test

import (
	"HRSystem/internal/model"
	"HRSystem/pkg/errors"
	"context"
	"reflect"
	"testing"
	"time"
)

func TestClockInRecordService_CreateClockInRecord(t *testing.T) {
	type args struct {
		ctx         context.Context
		accountID   uint64
		clockInType model.ClockInType
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Clock In",
			args: args{
				ctx:         context.Background(),
				accountID:   2,
				clockInType: model.ClockInTypeClockIn,
			},
			wantErr: nil,
		},
		{
			name: "Clock In without Clock Out",
			args: args{
				ctx:         context.Background(),
				accountID:   2,
				clockInType: model.ClockInTypeClockIn,
			},
			wantErr: errors.ErrAlreadyClockInToday,
		},
		{
			name: "Clock Out",
			args: args{
				ctx:         context.Background(),
				accountID:   2,
				clockInType: model.ClockInTypeClockOut,
			},
			wantErr: nil,
		},
		{
			name: "Clock Out without Clock In",
			args: args{
				ctx:         context.Background(),
				accountID:   3,
				clockInType: model.ClockInTypeClockOut,
			},
			wantErr: errors.ErrPleaseClockInFirst,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := clockInRecordSvc.CreateClockInRecord(tt.args.ctx, tt.args.accountID, tt.args.clockInType); err != tt.wantErr {
				t.Errorf("ClockInRecordService.CreateClockInRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClockInRecordService_ListClockInRecord(t *testing.T) {
	type args struct {
		ctx       context.Context
		accountID uint64
		start     time.Time
		end       time.Time
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantCount int
	}{
		{
			name: "First clock in, then list the clock-in records for the past hour, the count must be one",
			args: args{
				ctx:       context.Background(),
				accountID: 5,
				start:     time.Now().Add(-time.Hour),
				end:       time.Now().Add(time.Minute),
			},
			wantErr:   false,
			wantCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clockInRecordSvc.CreateClockInRecord(tt.args.ctx, tt.args.accountID, model.ClockInTypeClockIn)
			got, err := clockInRecordSvc.ListClockInRecord(tt.args.ctx, tt.args.accountID, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClockInRecordService.ListClockInRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.wantCount) {
				t.Errorf("ClockInRecordService.ListClockInRecord().Count = %v, wantCount %v", len(got), tt.wantCount)
			}
		})
	}
}
