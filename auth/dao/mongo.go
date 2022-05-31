package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mgo "server2/shared/mongo"
)

type Mongo struct {
	col      *mongo.Collection
	newObjID func() primitive.ObjectID
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col:      db.Collection("auth"),
		newObjID: primitive.NewObjectID,
	}
}

func (m *Mongo) ResolveAuthId(c context.Context, openId string) (string, error) {
	fmt.Printf("openId = %s \n", openId)
	_, err2 := m.col.InsertOne(c, bson.M{
		"_id":     m.newObjID(),
		"open_id": openId,
	})
	if err2 != nil {
		return "", err2
	}
	result := m.col.FindOneAndUpdate(c, bson.M{
		"open_id": openId,
	}, mgo.Set(bson.M{
		"open_id": openId,
	}), options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After),
	)
	if result.Err() != nil {
		return "", fmt.Errorf("FindOneAndUpdate err:%v", result.Err())
	}

	var row mgo.ObjId
	err := result.Decode(&row)
	if err != nil {
		return "", fmt.Errorf("解码失败 err:%v", err)
	}
	return row.ID.Hex(), nil
}
