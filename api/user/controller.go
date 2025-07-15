package user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web-scaffold/pkg/logger"
	"go-web-scaffold/pkg/models"
	"go-web-scaffold/pkg/service"
	"go-web-scaffold/pkg/utils"
	"strconv"
)

func Ping(c *gin.Context) {
	utils.SuccessWithData(c, "PONG!")
}

// /users?between=1&end=100
func GetByRange(c *gin.Context) {
	// 获取查询参数
	between, _ := strconv.Atoi(c.DefaultQuery("between", "1"))
	end, _ := strconv.Atoi(c.DefaultQuery("end", "100"))

	// 调用 service 层的方法来获取指定范围内的用户
	if err := service.GetGroupsOfUserByIdRange(int64(between), int64(end)); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c)
}

// Post，请求体带参数
func CreateUser(c *gin.Context) {
	var user models.User

	// 绑定 JSON 请求体
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	// 确保 name 和 email 不为空
	if user.Name == "" || user.Email == "" {
		err := errors.New("empty email or name")
		utils.Fail(c, err.Error())
		return
	}
	if err := service.CreateUser(user.Name, user.Email); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c)
}

// Put，请求体带参数
func UpdateUser(c *gin.Context) {
	var user models.User

	// 绑定 JSON 请求体
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	if user.Name == "" || user.Email == "" {
		err := errors.New("empty email or name")
		utils.Fail(c, err.Error())
		return
	}
	if err := service.UpdateUser(user.Name, user.Email); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c)
}

func DeleteUser(c *gin.Context) {
	fmt.Println(c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	if err = service.DeleteUserById(int64(id)); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c)
}

func SaveUser(c *gin.Context) {
	var user models.User

	// 绑定 JSON 请求体
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	if err = service.SaveUser(user.Name, user.Email, uint(id)); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c)
}

func RegisterRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	// 注册路由
	r.GET("/ping", Ping)
	r.GET("/user", GetByRange)               //查询参数
	r.POST("/user/create", CreateUser)       //请求体参数
	r.PUT("/user/update", UpdateUser)        //请求体参数
	r.DELETE("/user/delete/:id", DeleteUser) // 路径参数
	r.POST("/user/save/:id", SaveUser)       // 路径参数➕请求体参数
	return r
}
