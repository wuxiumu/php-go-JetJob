package api

import (
	"github.com/gin-gonic/gin"
	"jetjob/internal/storage"
	"net/http"
)

// 获取用户信息（需JWT）
func GetUserHandler(c *gin.Context) {
	email := c.GetString("email")
	var user User
	if err := storage.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// 修改用户信息（仅支持修改 name，可自行扩展）
type UpdateUserReq struct {
	Name string `json:"name" binding:"required"`
}

func UpdateUserHandler(c *gin.Context) {
	email := c.GetString("email")
	var req UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := storage.DB.Model(&User{}).Where("email = ?", email).
		Update("name", req.Name).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "修改失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "修改成功"})
}

// 删除用户
func DeleteUserHandler(c *gin.Context) {
	email := c.GetString("email")
	if err := storage.DB.Where("email = ?", email).Delete(&User{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
