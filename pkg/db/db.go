package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// Config contains the information required to connect to a database.
type Config struct {
	// Host specifies the URI to where the database is hosted.
	Host string

	// Username specifies the name of the database user to use when
	// establishing the connection to the database server.
	Username string

	// Password specifies the password of the database user used to establish
	// the connection to the database server.
	Password string

	// Name specifies the name of the database to use on the server.
	Name string

	// ConnectionTimeout specifies how many seconds to wait before giving up on
	// connecting to the database server.
	//
	// A timeout of zero means no timeout.
	ConnectionTimeout time.Duration
}

// DefaultConnectionTimeout represents the default amount of time to wait while
// establishing a connection to the database server.
const DefaultConnectionTimeout = 20 * time.Second

// Validate looks at the configuration's contents to ensure it has all the
// required fields.
func (conf *Config) validate() error {
	if conf.Host == "" {
		return errors.New("missing host")
	}

	if conf.Username == "" {
		return errors.New("missing username")
	}

	if conf.Password == "" {
		return errors.New("missing password")
	}

	if conf.Name == "" {
		return errors.New("missing name")
	}

	return nil
}

// DB represents a database. It contains a client used to connect to a database
// server and the database's collections.
type DB struct {
	client       *mongo.Client
	Reservations *mongo.Collection
}

const (
	reservationCollectionName = "reservations"
)

// New creates a database by establishing a connection to the database server
// specified in the given configuration.
func New(conf *Config) (*DB, error) {
	if conf == nil {
		return nil, errors.New("db: missing configuration")
	}

	err := conf.validate()
	if err != nil {
		return nil, fmt.Errorf("db: configuration %s", err)
	}

	url := fmt.Sprintf("mongodb://%s:%s@%s", conf.Username, conf.Password, conf.Host)
	client, err := mongo.NewClient(url)
	if err != nil {
		return nil, fmt.Errorf("db: failed to create client (%s)", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), conf.ConnectionTimeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("db: failed to connect to server (%s)", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("db: %s", err)
	}

	db := client.Database(conf.Name)
	if db == nil {
		return nil, fmt.Errorf("db: no database found with name \"%s\"", conf.Name)
	}

	reservations := db.Collection(reservationCollectionName)
	if reservations == nil {
		return nil, fmt.Errorf("db: no collection found with name \"%s\" in database", reservationCollectionName)
	}

	return &DB{client, reservations}, nil
}
