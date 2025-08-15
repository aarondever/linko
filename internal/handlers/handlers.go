package handlers

import (
	"github.com/aarondever/url-forg/internal/services"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	URLHandler *URLHandler
}

func InitializeHandlers(services *services.Services) *Handlers {
	// Initialize each handler - add new handlers here
	return &Handlers{
		URLHandler: NewURLHandler(services.URLService),
	}
}

func (handlers *Handlers) SetupRouters(router *chi.Mux) {
	// Setup API routes
	handlers.URLHandler.RegisterRoutes(router)
}
