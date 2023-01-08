package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client = mongo.Client
type Database = mongo.Database
type Collection = mongo.Collection

func Create(uri string) (*Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	return client, err
}
