package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func EnsureUserIndexes(
	ctx context.Context,
	db *mongodriver.Database,
) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	collection := db.Collection("users")

	indexes := []mongodriver.IndexModel{
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
			},
			Options: options.Index().
				SetUnique(true).
				SetName("unique_email"),
		},
		{
			Keys: bson.D{
				{Key: "username", Value: 1},
			},
			Options: options.Index().
				SetUnique(true).
				SetName("unique_username"),
		},
		{
			Keys: bson.D{
				{Key: "phone", Value: 1},
			},
			Options: options.Index().
				SetUnique(true).
				SetName("unique_phone"),
		},
		{
			Keys: bson.D{
				{Key: "alt_email", Value: 1},
			},
			Options: options.Index().
				SetUnique(true).
				SetPartialFilterExpression(
					bson.M{
						"alt_email": bson.M{
							"$gt": "",
						},
					},
				).
				SetName("unique_alt_email"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)

	return err
}
