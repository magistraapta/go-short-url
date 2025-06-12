package models

import "gorm.io/gorm"

type URL struct {
	gorm.Model
	ShortURL    string `gorm:"uniqueIndex;not null"`
	OriginalURL string `gorm:"not null"`
}
