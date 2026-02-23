package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
	"url_shorter/internal/entity"
	"url_shorter/internal/repo"
)

type Service struct {
	repo      *repo.Repository
	baseURL   string
	rateLimit int
}

func NewService(repo *repo.Repository, baseURL string, rateLimit int) *Service {
	return &Service{repo: repo, baseURL: baseURL, rateLimit: rateLimit}
}

func (s *Service) CheckRateLimit(ip string) (bool, error) {
	return s.repo.CheckRateLimit(ip, s.rateLimit)
}

func (s *Service) Shorten(originalURL string) (entity.CreateResponse, error) {
	code := generateRandomString(6)

	link := &entity.ShortLink{
		Code:      code,
		URL:       originalURL,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(link); err != nil {
		return entity.CreateResponse{}, err
	}

	return entity.CreateResponse{
		Code:     code,
		ShortURL: fmt.Sprintf("%s%s", s.baseURL, code),
	}, nil
}

func (s *Service) GetOriginalURL(code string) (string, error) {
	link, err := s.repo.GetByCode(code)
	if err != nil {
		return "", errors.New("ссылка не найдена")
	}

	go s.repo.IncrementClicks(code)

	return link.URL, nil
}

func (s *Service) GetStats(code string) (*entity.ShortLink, error) {
	return s.repo.GetByCode(code)
}

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
