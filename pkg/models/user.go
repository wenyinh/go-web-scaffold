package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `gorm:"size:64;not null;unique" json:"name"`
	Email string `gorm:"size:128;not null;unique" json:"email"`
}
