package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"azure.com/ecovo/reservation-service/cmd/handler"
	"azure.com/ecovo/reservation-service/cmd/middleware/auth"
	"azure.com/ecovo/reservation-service/pkg/db"
	"azure.com/ecovo/reservation-service/pkg/reservation"
	"azure.com/ecovo/reservation-service/pkg/trip"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	authConfig := auth.Config{
		Domain:               os.Getenv("AUTH_DOMAIN"),
		BasicAuthCredentials: os.Getenv("AUTH_CREDENTIALS"),
	}
	authBasicValidator, err := auth.NewBasicAuthValidator(&authConfig)
	if err != nil {
		log.Fatal(err)
	}
	authTokenValidator, err := auth.NewTokenValidator(&authConfig)
	if err != nil {
		log.Fatal(err)
	}
	authValidators := map[string]auth.Validator{
		"basic":  authBasicValidator,
		"bearer": authTokenValidator,
	}

	dbConnectionTimeout, err := time.ParseDuration(os.Getenv("DB_CONNECTION_TIMEOUT") + "s")
	if err != nil {
		dbConnectionTimeout = db.DefaultConnectionTimeout
	}
	dbConfig := db.Config{
		Host:              os.Getenv("DB_HOST"),
		Username:          os.Getenv("DB_USERNAME"),
		Password:          os.Getenv("DB_PASSWORD"),
		Name:              os.Getenv("DB_NAME"),
		ConnectionTimeout: dbConnectionTimeout}
	db, err := db.New(&dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	tripRepository, err := trip.NewRestRepository(os.Getenv("TRIP_SERVICE_DOMAIN"), os.Getenv("AUTH_CREDENTIALS"))
	if err != nil {
		log.Fatal(err)
	}

	tripUseCase := trip.NewService(tripRepository)

	reservationRepository, err := reservation.NewMongoRepository(db.Reservations)
	if err != nil {
		log.Fatal(err)
	}
	reservationUseCase := reservation.NewService(reservationRepository, tripUseCase)

	r := mux.NewRouter()

	// Reservations
	r.Handle("/reservations", handler.RequestID(handler.Auth(authValidators, handler.CreateReservation(reservationUseCase)))).
		Methods("POST").
		HeadersRegexp("Content-Type", "application/(json|json; charset=utf8)")
	r.Handle("/reservations/{id}", handler.RequestID(handler.Auth(authValidators, handler.DeleteReservation(reservationUseCase)))).
		Methods("DELETE").
		HeadersRegexp("Content-Type", "application/json")
	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r)))
}
