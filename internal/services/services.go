package services

import (
	"github.com/aarondever/linko/internal/config"
	"github.com/aarondever/linko/internal/database"
)

type Services struct {
	URLService *URLService
}

func InitializeServices(db *database.Database, cfg *config.Config) *Services {
	// Initialize each service - add new services here
	return &Services{
		URLService: NewURLService(db, cfg),
	}
}
