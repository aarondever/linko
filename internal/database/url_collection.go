package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/aarondever/url-forg/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log/slog"
	"os"
	"time"
)

const urlCollectionName = "urls"

func (database *Database) IsURLShortCodeExists(ctx context.Context, shortCode string) (bool, error) {
	var existing models.URLMapping
	if err := database.urlCollection.FindOne(ctx, bson.M{"short_code": shortCode}).Decode(&existing); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("error checking existing short code: %w", err)
	}
	return true, nil
}

func (database *Database) CreateURLShortCode(ctx context.Context, shortCode, url string) (*mongo.InsertOneResult, error) {
	mapping := models.URLMapping{
		ShortCode: shortCode,
		URL:       url,
		CreatedAt: time.Now(),
	}

	result, err := database.urlCollection.InsertOne(ctx, mapping)
	if err != nil {
		return nil, fmt.Errorf("error creating URL short code: %w", err)
	}
	return result, nil
}

func (database *Database) GetURL(ctx context.Context, shortCode string) (string, error) {
	var mapping models.URLMapping
	if err := database.urlCollection.FindOne(ctx, bson.M{"short_code": shortCode}).Decode(&mapping); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", fmt.Errorf("short code not found: %s", shortCode)
		}
		return "", fmt.Errorf("error retrieving URL mapping: %w", err)
	}
	return mapping.URL, nil
}

func (database *Database) initURLCollection(ctx context.Context) *mongo.Collection {
	if err := database.createCollection(ctx, urlCollectionName, bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"short_code", "url"},
			"properties": bson.M{
				"short_code": bson.M{
					"bsonType":    "string",
					"pattern":     "^[a-zA-Z0-9]{8}$",
					"description": "must be a string of exactly 8 alphanumeric characters",
				},
				"url": bson.M{
					"bsonType":    "string",
					"pattern":     "^https?://.+",
					"description": "must be a valid URL starting with http:// or https://",
				},
				"created_at": bson.M{
					"bsonType":    "date",
					"description": "timestamp when the URL was shortened",
				},
			},
		},
	}); err != nil {
		slog.Error("Failed to create URL collection", "error", err)
		os.Exit(1)
	}

	collection := database.db.Collection(urlCollectionName)

	if err := database.createIndexes(ctx, collection, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "short_code", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("short_code_unique"),
		},
		{
			Keys:    bson.D{{Key: "url", Value: 1}},
			Options: options.Index().SetName("url_index"),
		},
	}); err != nil {
		slog.Error("Failed to create URL indexes", "error", err)
		os.Exit(1)
	}

	return collection
}
