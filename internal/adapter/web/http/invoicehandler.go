package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shawnritchie/go-video-store/internal/port/driven"
	"io/ioutil"
	"net/http"
)

type (
	rental struct {
		Name string `json:"name"`
		Days uint16 `json:"days"`
	}

	returnRequest struct {
		Return []rental `json:"return"`
	}

	invoiceResponse struct {
		Return       []rental
		Price        uint64
		Currency     string
		MonetaryUnit string
	}
)

func (r *rental) isValid() bool {
	return !(r.Name == "" || r.Days <= 0)
}

func (r returnRequest) isValid() bool {
	for _, rental := range r.Return {
		if !rental.isValid() {
			return false
		}
	}

	return len(r.Return) > 0
}

func (s *server) processReturn(w http.ResponseWriter, r *http.Request) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("request body read error : %w", err)
	}

	var request returnRequest
	if err := json.Unmarshal(reqBody, &request); err != nil || !request.isValid() {
		return NewClientError(err, http.StatusBadRequest, "Bad Request: Post payload cannot be deserialized")
	}

	var returns []driven.FilmReturn
	for _, ele := range request.Return {
		returns = append(returns, driven.FilmReturn{FilmName: ele.Name, Days: ele.Days})
	}

	invoice, err := s.invoicer.Invoice(returns)
	if err != nil {
		switch {
		case errors.As(err, &driven.TypeInvalidRentalRequest):
			return NewClientError(err, http.StatusBadRequest, "Bad Request: submitted request cannot be processed!")
		default:
			return fmt.Errorf("error generating invoice: %w", err)
		}
	}

	setHeaders(w)
	json.NewEncoder(w).Encode(invoiceResponse{
		Return:       request.Return,
		Price:        uint64(invoice.Cost),
		Currency:     "SEK",
		MonetaryUnit: "Kr",
	})
	return nil
}
