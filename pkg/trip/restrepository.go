package trip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"azure.com/ecovo/reservation-service/pkg/entity"
)

// A RestRepository is a repository that performs HTTP requests on trips from the trip-service.
type RestRepository struct {
	domain    string
	authToken string
}

// NewRestRepository creates a REST repository.
func NewRestRepository(domain string, authToken string) (Repository, error) {
	if domain == "" {
		return nil, fmt.Errorf("trip.restrepository: domain is nil")
	}

	if authToken == "" {
		return nil, fmt.Errorf("trip.restrepository: authToken is nil")
	}

	return &RestRepository{domain, authToken}, nil
}

// CreateReservation creates a reservation on a trip.
func (r *RestRepository) CreateReservation(res *entity.Reservation) (*entity.Reservation, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(res)
	if err != nil {
		return nil, fmt.Errorf("trip.restrepository: failed to encode reservation")
	}

	req, err := http.NewRequest("POST", "http://"+r.domain+"/trips/"+res.TripID.Hex()+"/reservation", b)
	if err != nil {
		return nil, RequestError{fmt.Sprintf("trip.restrepository: failed to create request (%s)", err)}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", r.authToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("trip.restrepository: failed request to trip-service")
	}

	return res, nil
}

// DeleteReservation deletes a reservation from a trip.
func (r *RestRepository) DeleteReservation(res *entity.Reservation) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(res)
	if err != nil {
		return fmt.Errorf("trip.restrepository: failed to encode reservation")
	}

	req, err := http.NewRequest("DELETE", "https://"+r.domain+"/trips/{id}/reservation", b)
	if err != nil {
		return RequestError{fmt.Sprintf("trip.restrepository: failed to create request (%s)", err)}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", r.authToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("trip.restrepository: failed to validate token")
	}

	return nil
}
