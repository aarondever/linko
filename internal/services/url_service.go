package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/aarondever/url-forg/internal/config"
	"github.com/aarondever/url-forg/internal/database"
	"github.com/aarondever/url-forg/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"time"
)

const URLCollectionName = "urls"

type URLService struct {
	cfg           *config.Config
	db            *database.Database
	urlCollection *mongo.Collection
}

func NewURLService(cfg *config.Config, db *database.Database) *URLService {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.CreateCollectionWithValidation(ctx, URLCollectionName, models.GetURLValidator()); err != nil {
		return nil
	}

	collection := db.MongoDB.Collection(URLCollectionName)

	if err := db.CreateIndexes(ctx, collection, models.GetURLIndexes()); err != nil {
		return nil
	}

	return &URLService{
		cfg:           cfg,
		db:            db,
		urlCollection: collection,
	}
}

func (service *URLService) ShortenURL(ctx context.Context, url string) (string, error) {
	// Generate UUID and take first 8 characters
	uuidStr := uuid.New().String()
	shortCode := uuidStr[:8]

	// Check if short code already exists (handle collision)
	for {
		var existing models.URLMapping
		if err := service.urlCollection.FindOne(ctx, bson.M{"short_code": shortCode}).Decode(&existing); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				// Short code is available
				break
			}

			log.Printf("Error checking existing short code: %v", err)
			return "", err
		}

		// Collision detected, generate new short code
		uuidStr = uuid.New().String()
		shortCode = uuidStr[:8]
	}

	// Store in MongoDB
	mapping := models.URLMapping{
		ShortCode: shortCode,
		URL:       url,
		CreatedAt: time.Now(),
	}

	_, err := service.urlCollection.InsertOne(ctx, mapping)
	if err != nil {
		log.Printf("Error inserting URL mapping: %v", err)
		return "", err
	}

	return shortCode, nil
}

func (service *URLService) GetURL(ctx context.Context, shortCode string) (string, error) {
	var mapping models.URLMapping
	if err := service.urlCollection.FindOne(ctx, bson.M{"short_code": shortCode}).Decode(&mapping); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", fmt.Errorf("short code not found: %s", shortCode)
		}

		log.Printf("Error retrieving URL mapping: %v", err)
		return "", err
	}

	return mapping.URL, nil
}
