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
		utils.RespondWithError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	shortCode, err := handler.urlService.ShortenURL(request.Context(), params.URL)
	if err != nil {
		utils.RespondWithError(responseWriter, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(responseWriter, http.StatusOK, models.ShortenURLResponse{
		ShortCode: shortCode,
	})
}

func (handler *URLHandler) GetURL(responseWriter http.ResponseWriter, request *http.Request) {
	shortCode := request.PathValue("shortCode")
	originalURL, err := handler.urlService.GetURL(request.Context(), shortCode)
	if err != nil {
		utils.RespondWithError(responseWriter, http.StatusNotFound, "URL not found")
		return
	}

	utils.RespondWithJSON(responseWriter, http.StatusOK, models.GetURLResponse{
		OriginalURL: originalURL,
	})
}

func (handler *URLHandler) RedirectShortURL(responseWriter http.ResponseWriter, request *http.Request) {
	shortCode := request.PathValue("shortCode")
	originalURL, err := handler.urlService.GetURL(request.Context(), shortCode)
	if err != nil {
		utils.RespondWithError(responseWriter, http.StatusNotFound, "URL not found")
		return
	}

	http.Redirect(responseWriter, request, originalURL, http.StatusMovedPermanently)
}
