package repo

import (
	"errors"
	"time"
	"url_shorter/internal/entity"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(link *entity.ShortLink) error {
	return r.db.Create(link).Error
}

func (r *Repository) GetByCode(code string) (*entity.ShortLink, error) {
	var link entity.ShortLink
	err := r.db.Where("code = ?", code).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *Repository) IncrementClicks(code string) {
	r.db.Model(&entity.ShortLink{}).
		Where("code = ?", code).
		UpdateColumn("clicks", gorm.Expr("clicks + ?", 1))
}

func (r *Repository) CheckRateLimit(ip string, limit int) (bool, error) {
	var ipLimit entity.IpLimit

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("ip = ?", ip).First(&ipLimit).Error; err != nil {
			ipLimit = entity.IpLimit{
				IP:          ip,
				WindowStart: time.Now(),
				Count:       1,
			}
			return tx.Create(&ipLimit).Error
		}

		if time.Since(ipLimit.WindowStart) > time.Minute {
			ipLimit.Count = 1
			ipLimit.WindowStart = time.Now()
		} else {
			if ipLimit.Count >= limit {
				return errors.New("rate limit exceeded")
			}
			ipLimit.Count++
		}

		return tx.Save(&ipLimit).Error
	})

	if err != nil {
		if err.Error() == "rate limit exceeded" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
