package reservation

import (
	"context"
	"fmt"

	"azure.com/ecovo/reservation-service/pkg/entity"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

const (
	// DefaultRadius represents the default radius for a location search.
	DefaultRadius = 2000

	// TimeThreshold represents the time threshold for leaveAt or arriveBy (in hours)
	TimeThreshold = 12
)

// A MongoRepository is a repository that performs CRUD operations on reservations in
// a MongoDB collection.
type MongoRepository struct {
	collection *mongo.Collection
}

type document struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	TripID        primitive.ObjectID `bson:"tripId"`
	UserID        primitive.ObjectID `bson:"userId"`
	SourceID      primitive.ObjectID `bson:"sourceId"`
	DestinationID primitive.ObjectID `bson:"destinationId"`
	Seats         int                `bson:"seats"`
}

func newDocumentFromEntity(r *entity.Reservation) (*document, error) {
	if r == nil {
		return nil, fmt.Errorf("reservation.MongoRepository: entity is nil")
	}

	reservationID, err := getObjectID(r.ID)
	if err != nil {
		return nil, err
	}

	tripID, err := getObjectID(r.TripID)
	if err != nil {
		return nil, err
	}

	driverID, err := getObjectID(r.UserID)
	if err != nil {
		return nil, err
	}

	sourceID, err := getObjectID(r.SourceID)
	if err != nil {
		return nil, err
	}

	destinationID, err := getObjectID(r.DestinationID)
	if err != nil {
		return nil, err
	}

	return &document{
		reservationID,
		tripID,
		driverID,
		sourceID,
		destinationID,
		r.Seats,
	}, nil
}

func (d document) Entity() *entity.Reservation {
	return &entity.Reservation{
		entity.NewIDFromHex(d.ID.Hex()),
		entity.NewIDFromHex(d.TripID.Hex()),
		entity.NewIDFromHex(d.UserID.Hex()),
		entity.NewIDFromHex(d.SourceID.Hex()),
		entity.NewIDFromHex(d.DestinationID.Hex()),
		d.Seats,
	}
}

// NewMongoRepository creates a reservation repository for a MongoDB collection.
func NewMongoRepository(collection *mongo.Collection) (Repository, error) {
	if collection == nil {
		return nil, fmt.Errorf("reservation.MongoRepository: collection is nil")
	}

	return &MongoRepository{collection}, nil
}

// FindByID retrieves the reservation with the given ID, if it exists.
func (r *MongoRepository) FindByID(ID entity.ID) (*entity.Reservation, error) {
	objectID, err := primitive.ObjectIDFromHex(string(ID))
	if err != nil {
		return nil, fmt.Errorf("reservation.MongoRepository: failed to create object ID")
	}

	filter := bson.D{{"_id", objectID}}
	var d document
	err = r.collection.FindOne(context.TODO(), filter).Decode(&d)
	if err != nil {
		return nil, fmt.Errorf("reservation.MongoRepository: no reservation found with ID \"%s\" (%s)", ID, err)
	}

	return d.Entity(), nil
}

// Create stores the new reservation in the database and returns the unique
// identifier that was generated for it.
func (r *MongoRepository) Create(res *entity.Reservation) (entity.ID, error) {
	if res == nil {
		return entity.NilID, fmt.Errorf("reservation.MongoRepository: failed to create reservation (reservation is nil)")
	}

	d, err := newDocumentFromEntity(res)
	if err != nil {
		return entity.NilID, fmt.Errorf("reservation.MongoRepository: failed to create reservation document from entity (%s)", err)
	}

	resp, err := r.collection.InsertOne(context.TODO(), d)
	if err != nil {
		return entity.NilID, fmt.Errorf("reservation.MongoRepository: failed to create reservation (%s)", err)
	}

	ID, ok := resp.InsertedID.(primitive.ObjectID)
	if !ok {
		return entity.NilID, fmt.Errorf("reservation.MongoRepository: failed to get ID of created reservation (%s)", err)
	}

	return entity.ID(ID.Hex()), nil
}

// Update updates the reservation in the database.
func (r *MongoRepository) Update(res *entity.Reservation) error {
	d, err := newDocumentFromEntity(res)
	if err != nil {
		return fmt.Errorf("reservation.MongoRepository: failed to create reservation document from entity (%s)", err)
	}

	filter := bson.D{{"_id", d.ID}}
	update := bson.D{
		bson.E{"$set", d},
	}
	resp, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("reservation.MongoRepository: failed to update reservation with ID \"%s\" (%s)", res.ID, err)
	}

	if resp.MatchedCount <= 0 {
		return fmt.Errorf("reservation.MongoRepository: no matching reservation was found")
	}

	return nil
}

// Delete removes the reservation with the given ID from the database.
func (r *MongoRepository) Delete(ID entity.ID) error {
	objectID, err := primitive.ObjectIDFromHex(ID.Hex())
	if err != nil {
		return fmt.Errorf("reservation.MongoRepository: failed to create object ID")
	}

	filter := bson.D{{"_id", objectID}}
	_, err = r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("reservation.MongoRepository: failed to delete reservation with ID \"%s\" (%s)", ID, err)
	}

	return nil
}

// Gets an object ID from an entity of type ID
func getObjectID(rawID entity.ID) (primitive.ObjectID, error) {
	if rawID.IsZero() {
		return primitive.NilObjectID, nil
	}

	objectID, err := primitive.ObjectIDFromHex(rawID.Hex())
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("reservation.MongoRepository: failed to create object")
	}
	return objectID, nil
}
