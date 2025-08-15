package database

import (
	"context"
	"github.com/aarondever/url-forg/internal/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"time"
)

type Database struct {
	Mongo         *mongo.Client
	db            *mongo.Database
	urlCollection *mongo.Collection
}

func InitializeDatabase(config *config.Config) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(options.Client().ApplyURI(config.DatabaseURL))
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	// Test MongoDB connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		return nil, err
	}

	log.Println("Connected to MongoDB")

	database := &Database{
		Mongo: client,
		db:    client.Database(config.DBName),
	}

	// Initialize collections
	if err = database.initURLCollection(ctx); err != nil {
		return nil, err
	}

	return database, nil
}

func (database *Database) createCollection(ctx context.Context, collectionName string) error {
	// If collection exists, skip creation
	if database.isCollectionExists(ctx, collectionName) {
		return nil
	}

	// Create collection
	if err := database.db.CreateCollection(ctx, collectionName); err != nil {
		return err
	}

	log.Printf("Collection '%s' created successfully", collectionName)
	return nil
}

func (database *Database) createCollectionWithValidation(
	ctx context.Context,
	collectionName string,
	validator bson.M,
) error {
	// If collection exists, skip creation
	if database.isCollectionExists(ctx, collectionName) {
		return nil
	}

	// Create collection with validation schema
	opts := options.CreateCollection().SetValidator(validator)
	if err := database.db.CreateCollection(ctx, collectionName, opts); err != nil {
		return err
	}

	log.Printf("Collection '%s' created successfully", collectionName)
	return nil
}

func (database *Database) isCollectionExists(ctx context.Context, collectionName string) bool {
	collections, _ := database.db.ListCollectionNames(ctx, bson.M{"name": collectionName})

	return len(collections) > 0
}

func (database *Database) createIndexes(
	ctx context.Context,
	collection *mongo.Collection,
	indexes []mongo.IndexModel,
) error {
	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return err
	}

	return nil
}
