package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID       string `gorm:"unique"` // Change from int64 to string
	FirstName    string `json:"firstname" bson:"firstname"`
	LastName     string `json:"lastname" bson:"lastname"`
	UserName     string `json:"username" bson:"username" binding:"required"`
	LanguageCode string
}
