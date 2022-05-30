package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	col *mongo.Collection
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col: db.Collection("auth"),
	}
}

func (m *Mongo) ResolveAuthId(c context.Context, openId string) (string, error) {
	m.col.FindOneAndUpdate(c, bson.M{
		"open_id": openId,
	}, bson.M{
		"$set": bson.M{
			"open_id": openId,
		},
	}, options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After),
	)
	return "", nil
}
