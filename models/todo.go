package models

import "time"

// ShortURL represents a shortened URL record.
type ShortURL struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    OriginalURL string    `gorm:"not null" json:"original_url"`
    ShortCode   string    `gorm:"uniqueIndex;not null" json:"short_code"`
    Clicks      uint      `gorm:"default:0" json:"clicks"`
    CreatedAt   time.Time `json:"created_at"`
}
