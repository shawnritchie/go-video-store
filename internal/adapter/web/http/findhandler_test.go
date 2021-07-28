package http

import (
	"encoding/json"
	"fmt"
	"github.com/shawnritchie/go-video-store/internal/domain"
	"github.com/shawnritchie/go-video-store/internal/port/driven"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type spyFilmFinder struct {
	findInvocation uint32
	findParams     []string
	returnFx       func() (*domain.Film, error)
}

func (spy *spyFilmFinder) Find(name string) (*domain.Film, error) {
	spy.findInvocation++
	spy.findParams = append(spy.findParams, name)
	return spy.returnFx()
}

func newSpyFilmFinder(fx func() (*domain.Film, error)) *spyFilmFinder {
	return &spyFilmFinder{
		findInvocation: 0,
		findParams:     []string{},
		returnFx:       fx,
	}
}

const (
	FilmName     = "Loki"
	FilmDirector = "Marvel"
	FilmRelease  = domain.New
)

func TestFindRequest_Success(t *testing.T) {
	spy := newSpyFilmFinder(func() (*domain.Film, error) {
		film := &domain.Film{
			Name:     FilmName,
			Director: FilmDirector,
			Release:  FilmRelease,
		}
		return film, nil
	})
	server := New(spy, nil, nil)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/catalogue/film?name=%s", FilmName), nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler(server.findFilm)(res, req)

	var searchResponse findResponse
	unmarshalBody(res, &searchResponse)

	switch {
	case spy.findInvocation != 1:
		t.Errorf("was expecting single invocation but had %d invocations", spy.findInvocation)
	case spy.findParams[0] != FilmName:
		t.Errorf("was expecting the search to be executed on film %q", FilmName)
	case res.Code != http.StatusOK:
		t.Errorf("got status %d but wanted %d", res.Code, http.StatusOK)
	case searchResponse.Name != FilmName || searchResponse.Director != FilmDirector || searchResponse.Release != string(FilmRelease):
		t.Errorf("received unexpected response %#v", searchResponse)
	}
}

func TestFindRequest_MissingQueryParameter(t *testing.T) {
	spy := newSpyFilmFinder(func() (*domain.Film, error) {
		film := &domain.Film{
			Name:     FilmName,
			Director: FilmDirector,
			Release:  FilmRelease,
		}
		return film, nil
	})

	server := New(spy, nil, nil)

	req, err := http.NewRequest(http.MethodGet, "/catalogue/film", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	err = handler(server.findFilm)(res, req)

	clientError, ok := err.(ClientError)
	if !ok {
		t.Errorf("expected Client error but got %#v", err)
	}
	status, _ := clientError.ResponseHeaders()

	if status != http.StatusBadRequest {
		t.Errorf("got status %d but wanted %d", status, http.StatusBadRequest)
	}
}

func TestFindRequest_NoneCataloguedFilm(t *testing.T) {
	spy := newSpyFilmFinder(func() (*domain.Film, error) {
		return nil, &driven.FilmNotFoundError{Name: "Black Widow"}
	})

	server := New(spy, nil, nil)

	req, err := http.NewRequest(http.MethodGet, "/catalogue/film?name=Black%20Widow", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	err = handler(server.findFilm)(res, req)

	clientError, ok := err.(ClientError)
	if !ok {
		t.Errorf("expected Client error but got %#v", err)
	}
	status, _ := clientError.ResponseHeaders()

	if status != http.StatusNotFound {
		t.Errorf("got status %d but wanted %d", status, http.StatusNotFound)
	}
}

func unmarshalBody(w *httptest.ResponseRecorder, res interface{}) error {
	reqBody, err := ioutil.ReadAll(w.Body)
	if err != nil {
		return fmt.Errorf("request body read error : %w", err)
	}

	if err := json.Unmarshal(reqBody, &res); err != nil {
		return NewClientError(err, http.StatusBadRequest, "Bad Request: Post payload cannot be deserialized")
	}

	return nil
}
