package services

import (
	"context"
	"github.com/aarondever/linko/internal/config"
	"github.com/aarondever/linko/internal/database"
	"github.com/aarondever/linko/internal/models"
	"github.com/google/uuid"
)

type URLService struct {
	db  *database.Database
	cfg *config.Config
}

func NewURLService(db *database.Database, cfg *config.Config) *URLService {
	return &URLService{
		db:  db,
		cfg: cfg,
	}
}

func (service *URLService) ShortenURL(ctx context.Context, url string) (string, error) {
	// Generate UUID and take first 8 characters
	uuidStr := uuid.New().String()
	shortCode := uuidStr[:8]

	// Check if short code already exists (handle collision)
	for {
		exists, err := service.db.IsURLShortCodeExists(ctx, shortCode)
		if err != nil {
			return "", err
		}

		if !exists {
			// Short code is available
			break
		}

		// Collision detected, generate new short code
		uuidStr = uuid.New().String()
		shortCode = uuidStr[:8]
	}

	urlMapping := models.URLMapping{
		ShortCode: shortCode,
		URL:       url,
	}

	_, err := service.db.CreateURLShortCode(ctx, urlMapping)
	if err != nil {
		return "", err
	}

	return shortCode, nil
}

func (service *URLService) GetURL(ctx context.Context, shortCode string) (string, error) {
	url, err := service.db.GetURL(ctx, shortCode)
	if err != nil {
		return "", err
	}

	return url, nil
}
