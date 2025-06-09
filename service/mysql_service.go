package service

import (
	"errors"
	"fmt"
	"go-web-scaffold/dao/mysql"
	"go-web-scaffold/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateTableIfNotExists() (err error) {
	err = mysql.DB.AutoMigrate(&models.User{})
	if err != nil {
		zap.L().Error("create table failed", zap.Error(err))
		return
	}
	fmt.Println("create users table success (or already exist)")
	return
}

func CreateUser(name, email string) (err error) {
	user := models.User{Name: name, Email: email}
	err = mysql.DB.Create(&user).Error
	if err != nil {
		zap.L().Error("Insert user failed", zap.Error(err))
		return
	}
	fmt.Println("insert users success")
	return
}

func GetUserByName(name string) (*models.User, error) {
	var user models.User
	if err := mysql.DB.Model(&models.User{}).Where("name = ?", name).First(&user).Error; err != nil {
		zap.L().Error("get user failed", zap.Error(err))
		return nil, err
	}
	fmt.Println("get user success")
	return &user, nil
}

func GetGroupsOfUserByIdRange(between, end int64) (err error) {
	var users []models.User
	res := mysql.DB.Model(&models.User{}).Where("id between ? and ?", between, end).Find(&users)
	err = res.Error
	if err != nil {
		zap.L().Error("get users failed", zap.Error(res.Error))
		return res.Error
	}
	fmt.Println("get users success, rows affected:", res.RowsAffected)
	return
}

func UpdateUser(name, email string) (err error) {
	err = mysql.DB.Model(&models.User{}).Where("name = ?", name).Update("email", email).Error
	if err != nil {
		zap.L().Error("update user failed", zap.Error(err))
		return
	}
	fmt.Println("update user success")
	return
}

func SaveUser(name, email string, id uint) (err error) {
	if name == "" || email == "" {
		err = errors.New("name and email cannot be empty")
		return
	}
	user := models.User{Model: gorm.Model{ID: id}, Name: name, Email: email}
	res := mysql.DB.Model(&models.User{}).Save(&user)
	if res.Error != nil {
		zap.L().Error("save user failed", zap.Error(res.Error))
		return
	}
	fmt.Println("save user success")
	return
}

func DeleteUserById(id int64) (err error) {
	res := mysql.DB.Model(&models.User{}).Where("id = ?", id).Or("name = ?", "HWY").Delete(&models.User{})
	if res.Error != nil {
		zap.L().Error("delete user failed", zap.Error(res.Error))
		return
	}
	if res.RowsAffected == 0 {
		err = errors.New("delete user failed, no rows affected")
		zap.L().Error("delete user failed", zap.Error(err))
		return
	}
	fmt.Println("delete user success")
	return
}
