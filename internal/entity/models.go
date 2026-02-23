package entity

import "time"

type ShortLink struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	Code      string    `gorm:"uniqueIndex;type:varchar(10);not null" json:"code"`
	URL       string    `gorm:"type:text;not null" json:"url"`
	Clicks    int       `gorm:"default:0" json:"clicks"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateRequest struct {
	URL string `json:"url"`
}

type CreateResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"shortUrl"`
}

type IpLimit struct {
	IP          string    `gorm:"primaryKey;type:varchar(45)"`
	WindowStart time.Time `gorm:"type:timestamp"`
	Count       int       `gorm:"default:0"`
}
