package service

import (
	"errors"
	"fmt"
	"go-web-scaffold/dao/mysql"
	"go-web-scaffold/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateTableIfNotExists() error {
	return mysql.CreateUserTableIfNotExists()
}

func CreateUser(name, email string) error {
	user := models.User{Name: name, Email: email}
	if err := mysql.InsertUser(&user); err != nil {
		zap.L().Error("Insert user failed", zap.Error(err))
		return err
	}
	fmt.Println("insert users success")
	return nil
}

func GetUserByName(name string) (*models.User, error) {
	user, err := mysql.SelectUserByName(name)
	if err != nil {
		zap.L().Error("get user failed", zap.Error(err))
		return nil, err
	}
	fmt.Println("get user success")
	return user, nil
}

func GetGroupsOfUserByIdRange(between, end int64) error {
	users, err := mysql.SelectUsersByIdRange(between, end)
	if err != nil {
		zap.L().Error("get users failed", zap.Error(err))
		return err
	}
	fmt.Println("get users success, count:", len(users))
	return nil
}

func UpdateUser(name, email string) error {
	err := mysql.UpdateUserEmailByName(name, email)
	if err != nil {
		zap.L().Error("update user failed", zap.Error(err))
		return err
	}
	fmt.Println("update user success")
	return nil
}

func SaveUser(name, email string, id uint) error {
	if name == "" || email == "" {
		return errors.New("name and email cannot be empty")
	}
	user := models.User{Model: gorm.Model{ID: id}, Name: name, Email: email}
	if err := mysql.SaveUser(&user); err != nil {
		zap.L().Error("save user failed", zap.Error(err))
		return err
	}
	fmt.Println("save user success")
	return nil
}

func DeleteUserById(id int64) error {
	if err := mysql.DeleteUserById(id); err != nil {
		zap.L().Error("delete user failed", zap.Error(err))
		return err
	}
	fmt.Println("delete user success")
	return nil
}
