package handlers

import (
	"github.com/aarondever/url-forg/internal/models"
	"github.com/aarondever/url-forg/internal/services"
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
	shortenRequest := models.ShortenURLRequest{}
	if err := decodeRequestBody(request, &shortenRequest); err != nil {
		respondWithError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	shortCode, err := handler.urlService.ShortenURL(request.Context(), shortenRequest.URL)
	if err != nil {
		respondWithError(responseWriter, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(responseWriter, http.StatusOK, models.ShortenURLResponse{
		ShortCode: shortCode,
	})
}

func (handler *URLHandler) GetURL(responseWriter http.ResponseWriter, request *http.Request) {
	shortCode := request.PathValue("shortCode")
	originalURL, err := handler.urlService.GetURL(request.Context(), shortCode)
	if err != nil {
		respondWithError(responseWriter, http.StatusNotFound, "URL not found")
		return
	}

	respondWithJSON(responseWriter, http.StatusOK, models.GetURLResponse{
		OriginalURL: originalURL,
	})
}

func (handler *URLHandler) RedirectShortURL(responseWriter http.ResponseWriter, request *http.Request) {
	shortCode := request.PathValue("shortCode")
	originalURL, err := handler.urlService.GetURL(request.Context(), shortCode)
	if err != nil {
		respondWithError(responseWriter, http.StatusNotFound, "URL not found")
		return
	}

	http.Redirect(responseWriter, request, originalURL, http.StatusMovedPermanently)
}
