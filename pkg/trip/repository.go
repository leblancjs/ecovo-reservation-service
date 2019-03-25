package trip

import (
	"azure.com/ecovo/reservation-service/pkg/entity"
)

// Repository is an interface representing the ability to perform CRUD
// operations on trip-service.
type Repository interface {
	CreateReservation(res *entity.Reservation) (*entity.Reservation, error)
	DeleteReservation(res *entity.Reservation) error
}
