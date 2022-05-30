package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestMongo_ResolveAuthId(t *testing.T) {
	c := context.Background()
	connect, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://admin:123456@127.0.0.1:27017"))
	if err != nil {
		t.Fatalf("mongo连接错误 %v", err)
	}

	m := NewMongo(connect.Database("serverDemo"))
	id, err := m.ResolveAuthId(c, "123")
	if err != nil {
		t.Fatalf("数据库插入失败 %v", err)
	}

	want := "62946a36163e0a8f94c8a54a"
	if id != want {
		t.Fatalf("断言错误,测试不通过")
	}
}
