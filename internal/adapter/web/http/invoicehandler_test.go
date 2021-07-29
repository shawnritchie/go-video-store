package http

import (
	"fmt"
	"github.com/shawnritchie/go-video-store/internal/domain"
	"github.com/shawnritchie/go-video-store/internal/port/driven"
	"net/http"
	"net/http/httptest"
	"testing"
)

type spyFilmInvoicer struct {
	requests [][]driven.FilmReturn
	cost     domain.SEK
	err      error
}

func (s *spyFilmInvoicer) Invoice(request []driven.FilmReturn) (*domain.RentalInvoice, error) {
	s.requests = append(s.requests, request)

	var rentals []domain.Rental
	for _, film := range request {
		rentals = append(rentals, domain.Rental{
			Film: domain.Film{
				Name:     film.FilmName,
				Director: FilmDirector,
				Release:  domain.New,
			},
			Days: domain.Days(film.Days),
		})
	}

	return &domain.RentalInvoice{
		RentalReturn: domain.RentalReturn{
			Rentals: rentals,
		},
		Cost: s.cost,
	}, s.err
}

func NewSpyFilmInvoicer(cost domain.SEK, err error) *spyFilmInvoicer {
	return &spyFilmInvoicer{
		requests: [][]driven.FilmReturn{},
		cost:     cost,
		err:      err,
	}
}

func TestInvoicer_SuccessfullyProcessedReturn(t *testing.T) {
	totalCost := domain.SEK(20)
	spyInvoicer := NewSpyFilmInvoicer(totalCost, nil)
	server := New(nil, nil, spyInvoicer)

	returnReq := returnRequest{
		Return: []rental{
			{Name: "Loki", Days: 5},
			{Name: "Doctor Strange", Days: 3},
			{Name: "Jon Wick", Days: 1},
		},
	}
	req, err := http.NewRequest(http.MethodPost, "/store/return", toJSON(returnReq))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler(server.processReturn)(res, req)

	var invoiceRes invoiceResponse
	unmarshalBody(res, &invoiceRes)

	switch {
	case len(spyInvoicer.requests) != 1:
		t.Errorf("was expecting single invocation to append the film")
	case len(spyInvoicer.requests[0]) != len(returnReq.Return):
		t.Errorf("was expecting 3 films in the received request")
	case res.Code != http.StatusOK:
		t.Errorf("got status %d but wanted %d", res.Code, http.StatusOK)
	case len(invoiceRes.Return) != len(returnReq.Return):
		t.Errorf("received unexpected number of returns relative to the request %#v", invoiceRes.Return)
	case invoiceRes.MonetaryUnit != "Kr" || invoiceRes.Currency != "SEK" || invoiceRes.Price != uint64(totalCost):
		t.Errorf("received unexpected costings Currency: %v MonetaryUnit: %v, Total: %d",
			invoiceRes.Currency, invoiceRes.MonetaryUnit, invoiceRes.Price)
	}

	for i, rental := range invoiceRes.Return {
		if rental != returnReq.Return[i] {
			t.Errorf("was expecting rental response %#v but received %#v", returnReq.Return[i], rental)
		}
	}
}

func TestInvoicer_CorruptedRequestPayload(t *testing.T) {
	totalCost := domain.SEK(20)
	spyInvoicer := NewSpyFilmInvoicer(totalCost, nil)
	server := New(nil, nil, spyInvoicer)

	type corruptedRental struct {
		Title string `json:"title"`
		Hours uint16 `json:"hours"`
	}

	type corruptedRequest struct {
		Return []corruptedRental `json:"return"`
	}

	returnReq := corruptedRequest{
		Return: []corruptedRental{
			{Title: "Loki", Hours: 5},
			{Title: "Doctor Strange", Hours: 3},
			{Title: "Jon Wick", Hours: 1},
		},
	}

	req, err := http.NewRequest(http.MethodPost, "/store/return", toJSON(returnReq))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	err = server.processReturn(res, req)

	clientError, ok := err.(ClientError)
	if !ok {
		t.Errorf("expected Client error but got %#v", err)
	}
	status, _ := clientError.ResponseHeaders()

	if status != http.StatusBadRequest {
		t.Errorf("got status %d but wanted %d", status, http.StatusBadRequest)
	}
}

func TestInvoicer_InvalidRentalRequest(t *testing.T) {
	totalCost := domain.SEK(20)
	spyInvoicer := NewSpyFilmInvoicer(totalCost, &driven.InvalidRentalRequestError{fmt.Errorf("invalid film")})
	server := New(nil, nil, spyInvoicer)

	returnReq := returnRequest{
		Return: []rental{
			{Name: "Loki", Days: 5},
			{Name: "Doctor Strange", Days: 3},
			{Name: "Jon Wick", Days: 1},
		},
	}
	req, err := http.NewRequest(http.MethodPost, "/store/return", toJSON(returnReq))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	err = server.processReturn(res, req)

	clientError, ok := err.(ClientError)
	if !ok {
		t.Errorf("expected Client error but got %#v", err)
	}
	status, _ := clientError.ResponseHeaders()

	if status != http.StatusBadRequest {
		t.Errorf("got status %d but wanted %d", status, http.StatusBadRequest)
	}
}
