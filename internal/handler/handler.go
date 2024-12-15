package handler

import (
	"HRSystem/internal/model"
	"HRSystem/internal/service"
	"HRSystem/pkg/errors"
	"HRSystem/pkg/jwthelper"
	"HRSystem/pkg/middleware"
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
	v1Auth := v1.Group("", middleware.ValidateJWT())
	v1Auth.POST("/changePassword", h.ChangePassword)
	v1Auth.POST("/account", h.CreateAccount)
	v1Auth.POST("/clockInRecord", h.CreateClockInRecord)
	v1Auth.GET("/clockInRecord", h.ListClockInRecord)

	return h
}

func (h *Handler) Login(ctx *gin.Context) {
	var requestBody struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}
	token, err := h.accountSvc.Login(ctx.Request.Context(), requestBody.Account, requestBody.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if err == errors.ErrAccountOrPasswordIncorrect {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *Handler) ChangePassword(ctx *gin.Context) {
	claim, _ := jwthelper.GetClaim(ctx)
	var requestBody struct {
		OldPassword      string `json:"oldPassword"`
		NewPassword      string `json:"newPassword"`
		CheckNewPassword string `json:"checkNewPassword"`
	}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}
	err := h.accountSvc.ChangePassword(ctx.Request.Context(), claim.Account, requestBody.OldPassword, requestBody.NewPassword, requestBody.CheckNewPassword)
	if err != nil {
		status := http.StatusInternalServerError
		if err == errors.ErrAccountOrPasswordIncorrect {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) CreateAccount(ctx *gin.Context) {
	claim, _ := jwthelper.GetClaim(ctx)
	var requestBody struct {
		Account string `json:"account"`
		Name    string `json:"name"`
	}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	creator := h.accountSvc.ConvertClaimToAccount(claim)

	pwd, err := h.accountSvc.CreateAccount(ctx.Request.Context(), creator, requestBody.Name, requestBody.Account)
	if err != nil {
		status := http.StatusInternalServerError
		if err == errors.ErrAccountDuplicate {
			status = http.StatusBadRequest
		} else if err == errors.ErrInsufficientPrivilege {
			status = http.StatusForbidden
		}
		ctx.JSON(status, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"password": pwd,
	})

}

func (h *Handler) CreateClockInRecord(ctx *gin.Context) {
	claim, _ := jwthelper.GetClaim(ctx)
	var requestBody struct {
		Type model.ClockInType `json:"type"`
	}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	err := h.clockInRecordSvc.CreateClockInRecord(ctx.Request.Context(), claim.ID, requestBody.Type)
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
	claim, _ := jwthelper.GetClaim(ctx)
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

	records, err := h.clockInRecordSvc.ListClockInRecord(ctx.Request.Context(), claim.ID, startTime, endTime)
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
