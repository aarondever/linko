package handlers

import (
	"github.com/aarondever/linko/internal/models"
	"github.com/aarondever/linko/internal/services"
	"github.com/aarondever/linko/internal/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type URLHandler struct {
	urlService *services.URLService
}

func NewURLHandler(urlService *services.URLService) *URLHandler {
	return &URLHandler{urlService: urlService}
}

func (handler *URLHandler) RegisterRoutes(router *chi.Mux) {
	router.Route("/api/v1/url", func(router chi.Router) {
		router.Post("/shorten", handler.ShortenURL)
		router.Get("/shorten/{shortCode}", handler.GetURL)
	})

	router.Get("/r/{shortCode}", handler.RedirectShortURL)
}

func (handler *URLHandler) ShortenURL(responseWriter http.ResponseWriter, request *http.Request) {
	var params models.ShortenURLRequest
	if err := utils.DecodeRequestBody(request, &params); err != nil {
		utils.RespondWithError(responseWriter, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortCode, err := handler.urlService.ShortenURL(request.Context(), params.URL)
	if err != nil {
		utils.RespondWithError(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(responseWriter, models.ShortenURLResponse{
		ShortCode: shortCode,
	}, http.StatusCreated)
}

func (handler *URLHandler) GetURL(responseWriter http.ResponseWriter, request *http.Request) {
	shortCode := request.PathValue("shortCode")
	originalURL, err := handler.urlService.GetURL(request.Context(), shortCode)
	if err != nil {
		utils.RespondWithError(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if originalURL == "" {
		utils.RespondWithError(responseWriter, "URL not found", http.StatusNotFound)
		return
	}

	utils.RespondWithJSON(responseWriter, models.GetURLResponse{
		OriginalURL: originalURL,
	}, http.StatusOK)
}

func (handler *URLHandler) RedirectShortURL(responseWriter http.ResponseWriter, request *http.Request) {
	shortCode := request.PathValue("shortCode")
	originalURL, err := handler.urlService.GetURL(request.Context(), shortCode)
	if err != nil {
		utils.RespondWithError(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if originalURL == "" {
		utils.RespondWithError(responseWriter, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(responseWriter, request, originalURL, http.StatusMovedPermanently)
}
