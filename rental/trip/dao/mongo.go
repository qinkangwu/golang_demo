package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	rentalpb "server2/rental/api/gen/v1"
	mgo "server2/shared/mongo"
)

type Mongo struct {
	col *mongo.Collection
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col: db.Collection("trip"),
	}
}

type TripRecord struct {
	mgo.ObjId         `bson:"inline"`
	mgo.UpdateAtField `bson:"inline"`
	Trip              *rentalpb.Trip `bson:"trip"`
}

func (m *Mongo) CreateTrip(c context.Context, trip *rentalpb.Trip) (*TripRecord, error) {
	r := &TripRecord{
		Trip: trip,
	}
	r.ID = mgo.NewObjID()
	r.UpdateAt = mgo.UpdateAt()
	_, insertErr := m.col.InsertOne(c, r)
	if insertErr != nil {
		return nil, insertErr
	}
	return r, nil
}

func (m Mongo) GetTrip(c context.Context, tripId string, userId string) (*TripRecord, error) {
	findTripId, fromHexErr := primitive.ObjectIDFromHex(tripId)
	if fromHexErr != nil {
		return nil, fromHexErr
	}
	result := m.col.FindOne(c, bson.M{
		"_id":         findTripId,
		"trip.userid": userId,
	})
	if err := result.Err(); err != nil {
		return nil, err
	}
	var tr TripRecord
	err := result.Decode(&tr)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

func (m Mongo) GetTrips(c context.Context, s rentalpb.TripStatus, userId string) ([]*TripRecord, error) {
	filter := bson.M{
		"trip.userid": userId,
	}
	if s != rentalpb.TripStatus_TS_NOT_SPECIFIED {
		filter["trip.status"] = s
	}
	findAll, findAllErr := m.col.Find(c, filter)
	if findAllErr != nil {
		return nil, findAllErr
	}
	var trips []*TripRecord
	for findAll.Next(c) {
		var trip TripRecord
		err := findAll.Decode(&trip)
		if err != nil {
			return nil, err
		}
		trips = append(trips, &trip)
	}
	return trips, nil
}

func (m Mongo) UpdateTrip(c context.Context, tId string, userId string, updatedAt int64, trip *rentalpb.Trip) (err error) {
	objectID, fromIdErr := primitive.ObjectIDFromHex(tId)
	if fromIdErr != nil {
		return fromIdErr
	}
	_, updateErr := m.col.UpdateOne(c, bson.M{
		"_id":         objectID,
		"trip.userid": userId,
		"updateat":    updatedAt,
	}, mgo.Set(bson.M{
		"trip":     trip,
		"updateat": mgo.UpdateAt(),
	}))
	if updateErr != nil {
		return updateErr
	}
	return updateErr
}
