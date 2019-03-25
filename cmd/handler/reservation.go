package handler

import (
	"encoding/json"
	"net/http"

	"azure.com/ecovo/reservation-service/pkg/entity"
	"azure.com/ecovo/reservation-service/pkg/reservation"
	"github.com/gorilla/mux"
)

// CreateReservation handles a request to create a reservation.
func CreateReservation(service reservation.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		var res *entity.Reservation
		err := json.NewDecoder(r.Body).Decode(&res)
		if err != nil {
			return err
		}

		res, err = service.Register(res)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			_ = service.Delete(entity.ID(res.ID))

			return err
		}

		return nil
	}
}

// DeleteReservation handles a request to delete a reservation.
func DeleteReservation(service reservation.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)

		id := entity.NewIDFromHex(vars["id"])

		err := service.Delete(id)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)

		return nil
	}
}
