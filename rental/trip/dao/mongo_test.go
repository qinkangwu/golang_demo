package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	rentalpb "server2/rental/api/gen/v1"
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

func TestMongo_GetTrip(t *testing.T) {
	mongoUri = "mongodb://admin:123456@127.0.0.1:27017"
	c := context.Background()
	fmt.Println(mongoUri)
	connect, err := mongo.Connect(c, options.Client().ApplyURI(mongoUri))
	if err != nil {
		t.Fatalf("mongo连接错误 %v", err)
	}

	m := NewMongo(connect.Database("serverDemo"))
	fromHex, fromHexError := primitive.ObjectIDFromHex("629727eaa9d97dc27e00fe73")
	if fromHexError != nil {
		t.Fatalf("测试不通过 %v", fromHexError)
	}
	findTrip, findTripErr := m.GetTrip(c, fromHex.Hex(), "user1")
	if findTripErr != nil {
		t.Fatalf("测试不通过 %v", findTripErr)
	}
	t.Fatalf("测试不通过 %+v", findTrip)
}

func TestMongo_CreateTrip(t *testing.T) {
	mongoUri = "mongodb://admin:123456@127.0.0.1:27017"
	c := context.Background()
	fmt.Println(mongoUri)
	connect, err := mongo.Connect(c, options.Client().ApplyURI(mongoUri))
	if err != nil {
		t.Fatalf("mongo连接错误 %v", err)
	}

	m := NewMongo(connect.Database("serverDemo"))
	trip, createErr := m.CreateTrip(c, &rentalpb.Trip{
		UserId: "user1",
		CarId:  "car1",
		Start: &rentalpb.LocationStatus{
			Location: &rentalpb.Location{
				Latitude:  123,
				Longitude: 35,
			},
			FeeCent:      12,
			KmDriven:     12,
			LocationDesc: "测试1",
		},
		Current: &rentalpb.LocationStatus{
			Location: &rentalpb.Location{
				Latitude:  123,
				Longitude: 35,
			},
			FeeCent:      12,
			KmDriven:     12,
			LocationDesc: "测试1",
		},
		End: &rentalpb.LocationStatus{
			Location: &rentalpb.Location{
				Latitude:  123,
				Longitude: 35,
			},
			FeeCent:      12,
			KmDriven:     12,
			LocationDesc: "测试1",
		},
		Status: rentalpb.TripStatus_IN_PROGRESS,
	})
	if createErr != nil {
		t.Error("测试错误", createErr)
	}

	t.Errorf("%+v", trip)
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m, &mongoUri))
}
