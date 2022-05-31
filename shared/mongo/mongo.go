package mgo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Set(v interface{}) bson.M {
	return bson.M{
		"$set": v,
	}
}

func SetOnInsert(v interface{}) bson.M {
	return bson.M{
		"$setOnInsert": v,
	}
}

type ObjId struct {
	ID primitive.ObjectID `bson:"_id"`
}

const IdField = "_id"
