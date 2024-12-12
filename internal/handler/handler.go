package handler

import (
	"HRSystem/internal/model"
	"HRSystem/internal/service"
	"HRSystem/pkg/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	*gin.Engine
	accountSvc       *service.AccountService
	clockInRecordSvc *service.ClockInRecordService
}

func New(
	accountSvc *service.AccountService,
	clockInRecordSvc *service.ClockInRecordService,
) http.Handler {
	h := &Handler{
		Engine:           gin.New(),
		accountSvc:       accountSvc,
		clockInRecordSvc: clockInRecordSvc,
	}
	h.Use(gin.Recovery())

	api := h.Group("/api")
	v1 := api.Group("/v1")

	v1.POST("/login", h.Login)
	v1.POST("/changePassword", h.ChangePassword)
	v1.POST("/account", h.CreateAccount)
	v1.POST("/clockInRecord", h.CreateClockInRecord)
	v1.GET("/clockInRecord", h.ListClockInRecord)

	return h
}

func (h *Handler) Login(ctx *gin.Context) {
	panic("not implement")
}

func (h *Handler) ChangePassword(ctx *gin.Context) {
	panic("not implement")
}

func (h *Handler) CreateAccount(ctx *gin.Context) {
	panic("not implement")
}

func (h *Handler) CreateClockInRecord(ctx *gin.Context) {
	// TODO - 從jwt 取 accountID
	accountID := uint64(1)
	var requestBody struct {
		Type model.ClockInType `json:"type"`
	}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	err := h.clockInRecordSvc.CreateClockInRecord(ctx.Request.Context(), accountID, requestBody.Type)
	if err != nil {
		status := http.StatusInternalServerError
		if err == errors.ErrAlreadyClockInToday || err == errors.ErrPleaseClockInFirst {
			status = http.StatusBadRequest
		}

		ctx.JSON(status, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{})
}

func (h *Handler) ListClockInRecord(ctx *gin.Context) {
	// TODO - 從jwt 取 accountID
	accountID := uint64(1)
	start := ctx.Query("start")
	end := ctx.Query("end")
	startTime, err := time.Parse(time.DateTime, start)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	endTime, err := time.Parse(time.DateTime, end)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	records, err := h.clockInRecordSvc.ListClockInRecord(ctx.Request.Context(), accountID, startTime, endTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"records": records,
	})
}
