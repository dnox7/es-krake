package mdb

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Sort(d bson.D) {
	opts := options.Find().SetSort(d)
}
