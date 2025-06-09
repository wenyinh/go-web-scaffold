package verify

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go-web-scaffold/service"
)

func TestMySQL() error {
	if err := service.CreateTableIfNotExists(); err != nil {
		return err
	}
	uid := uuid.New().String()
	name := "user_" + uid[:8]
	email := uid[:8] + "@verify.com"
	if err := service.CreateUser(name, email); err != nil {
		return err
	}
	user, err := service.GetUserByName(name)
	if err != nil {
		return err
	}
	if user.Name == name && user.Email == email {
		fmt.Printf("verify mysql success: %+v\n", user)
		return nil
	} else {
		fmt.Println("verify fail, mismatch")
		return errors.New("verify fail, mismatch")
	}
}
