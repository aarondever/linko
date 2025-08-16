package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
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
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	ShortCode string        `bson:"short_code" json:"short_code"`
	URL       string        `bson:"url" json:"url"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}
