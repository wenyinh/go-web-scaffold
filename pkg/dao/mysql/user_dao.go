package mysql

import (
	"errors"
	"go-web-scaffold/pkg/models"
)

func CreateUserTableIfNotExists() error {
	return db.AutoMigrate(&models.User{})
}

func InsertUser(m *models.User) error {
	return db.Create(m).Error
}

func SelectUserByName(name string) (*models.User, error) {
	var user models.User
	err := db.Where("name = ?", name).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SelectUsersByIdRange(between int64, end int64) ([]models.User, error) {
	var users []models.User
	err := db.Where("id BETWEEN ? AND ?", between, end).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func UpdateUserEmailByName(name string, email string) error {
	return db.Model(&models.User{}).Where("name = ?", name).Update("email", email).Error
}

func SaveUser(m *models.User) error {
	if m.Name == "" || m.Email == "" {
		return errors.New("name or email cannot be empty")
	}
	return db.Save(m).Error
}

func DeleteUserById(id int64) error {
	res := db.Where("id = ?", id).Delete(&models.User{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("delete user failed, no rows affected")
	}
	return nil
}
