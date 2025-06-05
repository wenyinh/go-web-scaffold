package verify

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go-web-scaffold/dao/mysql"
	"go-web-scaffold/models"
	"go.uber.org/zap"
)

func createTableIfNotExists() (err error) {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(64) NOT NULL,
		email VARCHAR(128) NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = mysql.DB.Exec(createTableSQL)
	if err != nil {
		zap.L().Error("create table failed", zap.Error(err))
		return
	}
	fmt.Println("create users table success (or already exist)")
	return
}

func createUser(name, email string) (err error) {
	if err = createTableIfNotExists(); err != nil {
		zap.L().Error("create table failed", zap.Error(err))
		return
	}
	user := models.User{Name: name, Email: email}
	_, err = mysql.DB.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		zap.L().Error("Insert user failed", zap.Error(err))
		return
	}
	fmt.Println("insert users success")
	return
}

func getUserByName(name string) (*models.User, error) {
	var user models.User
	if err := mysql.DB.Get(&user, "SELECT * FROM users WHERE name = ?", name); err != nil {
		zap.L().Error("get user failed", zap.Error(err))
		return nil, err
	}
	fmt.Println("get user success")
	return &user, nil
}

func TestMySQL() error {
	uid := uuid.New().String()
	name := "user_" + uid[:8]
	email := uid[:8] + "@test.com"
	if err := createUser(name, email); err != nil {
		return err
	}
	user, err := getUserByName(name)
	if err != nil {
		return err
	}
	if user.Name == name && user.Email == email {
		fmt.Printf("test mysql success: %+v\n", user)
		return nil
	} else {
		fmt.Println("test fail, mismatch")
		return errors.New("test fail, mismatch")
	}
}
