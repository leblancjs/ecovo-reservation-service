package reservation

import (
	"azure.com/ecovo/reservation-service/pkg/entity"
)

// Repository is an interface representing the ability to perform CRUD
// operations on reservations in a database.
type Repository interface {
	FindByID(ID entity.ID) (*entity.Reservation, error)
	Create(reservation *entity.Reservation) (entity.ID, error)
	Update(reservation *entity.Reservation) error
	Delete(ID entity.ID) error
}
