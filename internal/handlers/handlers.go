package handlers

import "github.com/aarondever/url-forg/internal/services"

type Handlers struct {
	URLHandler *URLHandler
}

func InitializeHandlers(services *services.Services) *Handlers {
	// Initialize each handler - add new handlers here
	return &Handlers{
		URLHandler: NewURLHandler(services.URLService),
	}
}
