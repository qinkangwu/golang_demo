package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	mongotesting "server2/shared/testing"
	"testing"
)

var mongoUri string

func TestMongo_ResolveAuthId(t *testing.T) {
	c := context.Background()
	fmt.Println(mongoUri)
	connect, err := mongo.Connect(c, options.Client().ApplyURI(mongoUri))
	if err != nil {
		t.Fatalf("mongo连接错误 %v", err)
	}

	m := NewMongo(connect.Database("serverDemo"))
	m.newObjID = func() primitive.ObjectID {
		hex, err := primitive.ObjectIDFromHex("62946a36163e0a8f94c8a54a")
		if err != nil {
			return [12]byte{}
		}
		return hex
	}
	id, err := m.ResolveAuthId(c, "123")
	if err != nil {
		t.Fatalf("数据库插入失败 %v", err)
	}

	want := m.newObjID
	if id != want {
		t.Fatalf("断言错误,测试不通过")
	}
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m, &mongoUri))
}
