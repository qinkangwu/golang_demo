package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	mgo "server2/shared/mongo"
	mongotesting "server2/shared/testing"
	"testing"
)

var mongoUri string

func mustObjId(s string) primitive.ObjectID {
	hex, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return [12]byte{}
	}
	return hex
}

func TestMongo_ResolveAuthId(t *testing.T) {
	c := context.Background()
	fmt.Println(mongoUri)
	connect, err := mongo.Connect(c, options.Client().ApplyURI(mongoUri))
	if err != nil {
		t.Fatalf("mongo连接错误 %v", err)
	}

	m := NewMongo(connect.Database("serverDemo"))
	_, err = m.col.InsertMany(c, []interface{}{
		bson.M{
			"_id":     mustObjId("62946a36163e0a8f94c8a54b"),
			"open_id": "open_id_1",
		},
		bson.M{
			"_id":     mustObjId("62946a36163e0a8f94c8a54c"),
			"open_id": "open_id_2",
		},
	})
	if err != nil {
		panic(err)
	}
	mgo.NewObjID = func() primitive.ObjectID {
		hex := mustObjId("62946a36163e0a8f94c8a54a")
		return hex
	}
	cases := []struct {
		name   string
		openId string
		want   string
	}{
		{
			name:   "existing_user",
			openId: "open_id_1",
			want:   "62946a36163e0a8f94c8a54b",
		},
		{
			name:   "existing_user2",
			openId: "open_id_2",
			want:   "62946a36163e0a8f94c8a54c",
		},
		{
			name:   "no_existing_user",
			openId: "open_id_3",
			want:   "62946a36163e0a8f94c8a54a",
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			id, err := m.ResolveAuthId(context.Background(), cc.openId)
			if err != nil {
				t.Fatalf("数据库插入失败 %v", err)
			}
			want := cc.want
			if id != want {
				t.Fatalf("断言错误,测试不通过")
			}
		})
	}
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m, &mongoUri))
}
