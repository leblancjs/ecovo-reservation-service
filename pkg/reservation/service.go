package reservation

import (
	"fmt"

	"azure.com/ecovo/reservation-service/pkg/entity"
	"azure.com/ecovo/reservation-service/pkg/trip"
)

// UseCase is an interface representing the ability to handle the business
// logic that involves reservations.
type UseCase interface {
	Register(r *entity.Reservation) (*entity.Reservation, error)
	FindByID(ID entity.ID) (*entity.Reservation, error)
	Delete(ID entity.ID) error
}

// A Service handles the business logic related to reservations.
type Service struct {
	repo        Repository
	tripService trip.UseCase
}

// NewService creates a reservation service to handle business logic and manipulate
// reservations through a repository.
func NewService(repo Repository, tripService trip.UseCase) *Service {
	return &Service{repo, tripService}
}

// Register modifies reservation repository based on a reservation done.
func (s *Service) Register(r *entity.Reservation) (*entity.Reservation, error) {
	if r == nil {
		return nil, fmt.Errorf("reservation.Service: reservation is nil")
	}

	_, err := s.FindByID(r.ID)
	if err == nil {
		return nil, AlreadyExistsError{fmt.Sprintf("reservation.Service: reservation already exists with ID \"%s\"", r.ID)}
	}

	err = r.Validate()
	if err != nil {
		return nil, err
	}

	r.ID, err = s.repo.Create(r)
	if err != nil {
		return nil, err
	}

	err = s.tripService.RegisterReservation(r)
	if err != nil {
		delErr := s.repo.Delete(r.ID)
		if delErr != nil {
			return nil, delErr
		}
		return nil, err
	}

	return r, nil
}

// FindByID retrieves the reservation with the given ID in the repository, if it
// exists.
func (s *Service) FindByID(ID entity.ID) (*entity.Reservation, error) {
	r, err := s.repo.FindByID(ID)
	if err != nil {
		return nil, NotFoundError{err.Error()}
	}

	return r, nil
}

// Delete modifies reservation respository based on a reservation done.
func (s *Service) Delete(ID entity.ID) error {
	res, err := s.repo.FindByID(ID)
	if err != nil {
		return err
	}

	err = s.tripService.DeleteReservation(res)
	if err != nil {
		return err
	}

	err = s.repo.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}
