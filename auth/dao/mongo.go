package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mgo "server2/shared/mongo"
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
	fmt.Printf("openId = %s \n", openId)
	result := m.col.FindOneAndUpdate(c, bson.M{
		"open_id": openId,
	}, mgo.SetOnInsert(bson.M{
		"open_id": openId,
		"_id":     mgo.NewObjID(),
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
