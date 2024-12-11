package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	*gin.Engine
}

func New() http.Handler {
	h := &Handler{
		Engine: gin.New(),
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
	panic("not implement")
}

func (h *Handler) ListClockInRecord(ctx *gin.Context) {
	panic("not implement")
}
