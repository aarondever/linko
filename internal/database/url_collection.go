package database

import (
	"context"
	"errors"
	"github.com/aarondever/linko/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log/slog"
	"time"
)

const urlCollectionName = "urls"

func (database *Database) IsURLShortCodeExists(ctx context.Context, shortCode string) (bool, error) {
	var existing models.URLMapping
	if err := database.urlCollection.FindOne(ctx, bson.M{"short_code": shortCode}).Decode(&existing); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}

		slog.Error("Failed checking existing short code", "error", err)
		return false, err
	}

	return true, nil
}

func (database *Database) GetURLMappingByID(ctx context.Context, id string) (*models.URLMapping, error) {
	mappingID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		slog.Error("Failed parse mapping ID", "error", err)
		return nil, err
	}

	var mapping models.URLMapping
	if err = database.urlCollection.FindOne(ctx, bson.M{"_id": mappingID}).Decode(&mapping); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			slog.Error("Failed find URL mapping", "error", err)
		}

		return nil, err
	}

	return &mapping, nil
}

func (database *Database) CreateURLShortCode(
	ctx context.Context,
	params models.URLMapping,
) (*models.URLMapping, error) {
	params.CreatedAt = time.Now()

	result, err := database.urlCollection.InsertOne(ctx, params)
	if err != nil {
		slog.Error("Failed insert URL short code", "error", err)
		return nil, err
	}

	return database.GetURLMappingByID(ctx, result.InsertedID.(bson.ObjectID).Hex())
}

func (database *Database) GetURL(ctx context.Context, shortCode string) (string, error) {
	var mapping models.URLMapping
	if err := database.urlCollection.FindOne(ctx, bson.M{"short_code": shortCode}).Decode(&mapping); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			slog.Error("Failed find URL mapping", "error", err)
		}

		return "", err
	}

	return mapping.URL, nil
}

func (database *Database) initURLCollection(ctx context.Context) *mongo.Collection {
	database.createCollection(ctx, urlCollectionName, bson.M{
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
	})

	collection := database.db.Collection(urlCollectionName)

	database.createIndexes(ctx, collection, []mongo.IndexModel{
		// Index on created_at for chronological queries
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("created_at_desc"),
		},
		// Index on short_code for finding url by code
		{
			Keys:    bson.D{{Key: "short_code", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("short_code_unique"),
		},
	})

	return collection
}
