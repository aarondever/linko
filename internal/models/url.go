package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

type ShortenURLRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type ShortenURLResponse struct {
	ShortCode string `json:"short_code"`
}

type GetURLResponse struct {
	OriginalURL string `json:"original_url"`
}

// URLMapping represents the URL document in MongoDB
type URLMapping struct {
	ShortCode string    `bson:"short_code" json:"short_code"`
	URL       string    `bson:"url" json:"url"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

func GetURLValidator() bson.M {
	return bson.M{
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
	}
}

func GetURLIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "short_code", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("short_code_unique"),
		},
		{
			Keys:    bson.D{{Key: "url", Value: 1}},
			Options: options.Index().SetName("url_index"),
		},
	}
}
