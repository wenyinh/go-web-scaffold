package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web-scaffold/logger"
	"go-web-scaffold/models"
	"go-web-scaffold/service"
	"net/http"
	"strconv"
)

func ping(c *gin.Context) {
	c.String(http.StatusOK, "PONG!")
}

// /users?between=1&end=100
func getByRange(c *gin.Context) {
	// 获取查询参数
	between, _ := strconv.Atoi(c.DefaultQuery("between", "1"))
	end, _ := strconv.Atoi(c.DefaultQuery("end", "100"))

	// 调用 service 层的方法来获取指定范围内的用户
	if err := service.GetGroupsOfUserByIdRange(int64(between), int64(end)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Users fetched successfully"})
}

// Post，请求体带参数
func createUser(c *gin.Context) {
	var user models.User

	// 绑定 JSON 请求体
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 确保 name 和 email 不为空
	if user.Name == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and email are required"})
		return
	}
	if err := service.CreateUser(user.Name, user.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// Put，请求体带参数
func updateUser(c *gin.Context) {
	var user models.User

	// 绑定 JSON 请求体
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if user.Name == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and email are required"})
		return
	}
	if err := service.UpdateUser(user.Name, user.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func deleteUser(c *gin.Context) {
	fmt.Println(c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = service.DeleteUserById(int64(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func saveUser(c *gin.Context) {
	var user models.User

	// 绑定 JSON 请求体
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = service.SaveUser(user.Name, user.Email, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User saved successfully"})
}

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	// 注册路由
	r.GET("/ping", ping)
	r.GET("/user", getByRange)               //查询参数
	r.POST("/user/create", createUser)       //请求体参数
	r.PUT("/user/update", updateUser)        //请求体参数
	r.DELETE("/user/delete/:id", deleteUser) // 路径参数
	r.POST("/user/save/:id", saveUser)       // 路径参数➕表单参数                            //请求体参数+路径参数
	return r
}
