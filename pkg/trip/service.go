package trip

import (
	"fmt"

	"azure.com/ecovo/reservation-service/pkg/entity"
)

// UseCase is an interface representing the ability to handle the business
// logic that involves trips.
type UseCase interface {
	RegisterReservation(r *entity.Reservation) error
	DeleteReservation(r *entity.Reservation) error
}

// A Service handles the business logic related to trips.
type Service struct {
	repo Repository
}

// NewService creates a trip service to handle business logic and manipulate
// trips through a repository.
func NewService(repo Repository) *Service {
	return &Service{repo}
}

// RegisterReservation will send a creation request to the rest repository that communicates with the trip-service.
func (s *Service) RegisterReservation(r *entity.Reservation) error {
	if r == nil {
		return fmt.Errorf("trip.Service: reservation is nil")
	}

	res, err := s.repo.CreateReservation(r)
	if err != nil {
		return err
	}

	err = res.Validate()
	if err != nil {
		return err
	}

	return nil
}

// DeleteReservation will send a deletion request to the rest repository that communicates with the trip-service.
func (s *Service) DeleteReservation(r *entity.Reservation) error {
	if r == nil {
		return fmt.Errorf("trip.Service: reservation is nil")
	}

	err := s.repo.DeleteReservation(r)
	if err != nil {
		return err
	}

	return nil
}
