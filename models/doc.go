package models

import (
	"time"

	"gorm.io/gorm"
)

type DocumentEmbedding struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Filename  string
	DocID     string
	DocText   string
	//Embedding []float64 `gorm:"type:float8[]"`
	Embedding string `gorm:"type:float8[]"` // Store as a string and ensure it's passed correctly
}
