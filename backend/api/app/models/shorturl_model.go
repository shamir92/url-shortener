package models

import (
	"time"

	"github.com/google/uuid"
)

type ShortUrl struct {
	ID        uuid.UUID `gorm:"primaryKey"`                                      // ID of ShortUrl.
	LongUrl   string    `validate:"required" gorm:"not null" json:"long-url"`    // LongUrl of ShortUrl.
	ShortUrl  string    `gorm:"uniqueIndex;not null"`                            // ShortUrl of ShortUrl.
	Email     string    `validate:"required,email" gorm:"not null" json:"email"` // User's email.
	CreatedAt time.Time // CreatedAt of ShortUrl.
	UpdatedAt time.Time // UpdatedAt of ShortUrl.
	ExpiredAt time.Time // ExpiredAt of ShortUrl.
}
