package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dysodeng/app/internal/application/service"
	"github.com/dysodeng/app/internal/interfaces/http/dto/response/api"
)

// UserHandler 用户HTTP处理器
type UserHandler struct {
	userAppService *service.UserAppService
}

// NewUserHandler 创建用户HTTP处理器
func NewUserHandler(userAppService *service.UserAppService) *UserHandler {
	return &UserHandler{
		userAppService: userAppService,
	}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userAppService.RegisterUser(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "用户注册成功",
		"user_id": user.ID,
	})
}

// GetUser 获取用户信息
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	userId, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusOK, api.Fail(c, "无效的用户ID", api.CodeFail))
		return
	}
	user, err := h.userAppService.GetUser(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

// ListUsers 获取用户列表
func (h *UserHandler) ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := h.userAppService.GetUserList(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []gin.H
	for _, user := range users {
		result = append(result, gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
		"data":      result,
	})
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	userId, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusOK, api.Fail(c, "无效的用户ID", api.CodeFail))
		return
	}

	err = h.userAppService.DeleteUser(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}
