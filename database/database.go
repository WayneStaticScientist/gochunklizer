package database

import "go.mongodb.org/mongo-driver/mongo"

type Database struct {
	db mongo.Client
}
