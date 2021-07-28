package http

import (
	"github.com/shawnritchie/go-video-store/internal/port/driven"
	"net/http"
	"net/http/httptest"
	"testing"
)

type spyFilmAppender struct {
	invocations []struct {
		name     string
		director string
	}
	throw error
}

func (s *spyFilmAppender) invoke(name string, director string) error {
	s.invocations = append(s.invocations, struct {
		name     string
		director string
	}{name: name, director: director})
	return s.throw
}

func (s *spyFilmAppender) AddNew(name string, director string) error {
	return s.invoke(name, director)
}

func (s *spyFilmAppender) AddRegular(name string, director string) error {
	return s.invoke(name, director)
}

func (s *spyFilmAppender) AddOld(name string, director string) error {
	return s.invoke(name, director)
}

func newSpyFilmAppender(throw error) *spyFilmAppender {
	return &spyFilmAppender{
		invocations: []struct {
			name     string
			director string
		}{},
		throw: throw,
	}
}

func TestFilmAppeneder_Success(t *testing.T) {
	spyAppender := newSpyFilmAppender(nil)
	server := New(nil, spyAppender, nil)

	appendReq := appendRequest{Name: FilmName, Director: FilmDirector}
	req, err := http.NewRequest(http.MethodPost, "catalogue/film/regular", toJSON(appendReq))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler(server.addRegularFilm)(res, req)

	var appendResponse appendResponse
	unmarshalBody(res, &appendResponse)

	switch {
	case len(spyAppender.invocations) != 1:
		t.Errorf("was expecting single invocation to append the film")
	case spyAppender.invocations[0].name != FilmName || spyAppender.invocations[0].director != FilmDirector:
		t.Errorf("was expecting the search to be executed on film %#v", appendReq)
	case res.Code != http.StatusOK:
		t.Errorf("got status %d but wanted %d", res.Code, http.StatusOK)
	case appendResponse.Name != FilmName || appendResponse.Director != FilmDirector || appendResponse.Release != "Regular":
		t.Errorf("received unexpected response %#v", appendResponse)
	}
}

func TestFilmAppeneder_CorruptedPayload(t *testing.T) {
	spyAppender := newSpyFilmAppender(nil)
	server := New(nil, spyAppender, nil)

	appendReq := struct {
		FilmName string `json:"filmName"`
		Director string `json:"dir"`
	}{
		FilmName: FilmName,
		Director: FilmDirector,
	}

	req, err := http.NewRequest(http.MethodPost, "catalogue/film/regular", toJSON(appendReq))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	err = server.addRegularFilm(res, req)

	clientError, ok := err.(ClientError)
	if !ok {
		t.Errorf("expected Client error but got %#v", err)
	}
	status, _ := clientError.ResponseHeaders()

	if status != http.StatusBadRequest {
		t.Errorf("got status %d but wanted %d", status, http.StatusBadRequest)
	}
}

func TestFilmAppeneder_AlreadyExists(t *testing.T) {
	spyAppender := newSpyFilmAppender(&driven.FilmAlreadyExist{Name: FilmName})
	server := New(nil, spyAppender, nil)

	appendReq := appendRequest{Name: FilmName, Director: FilmDirector}
	req, err := http.NewRequest(http.MethodPost, "catalogue/film/regular", toJSON(appendReq))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	err = server.addRegularFilm(res, req)

	clientError, ok := err.(ClientError)
	if !ok {
		t.Errorf("expected Client error but got %#v", err)
	}
	status, _ := clientError.ResponseHeaders()

	if status != http.StatusConflict {
		t.Errorf("got status %d but wanted %d", status, http.StatusConflict)
	}
}
