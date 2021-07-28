package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shawnritchie/go-video-store/internal/domain"
	"github.com/shawnritchie/go-video-store/internal/port/driven"
	"io/ioutil"
	"net/http"
)

type (
	appendRequest struct {
		Name     string `json:"name"`
		Director string `json:"director"`
	}

	appendResponse struct {
		Name     string `json:"name"`
		Director string `json:"director"`
		Release  string `json:"release"`
	}
)

func (ap *appendRequest) validate() bool {
	return ap.Name == "" || ap.Director == ""
}

func (s *server) addNewFilm(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if ct != contentType {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//Handle Corrupted Payload
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var film appendRequest
	if err := json.Unmarshal(reqBody, &film); err != nil {
		//Handle Corrupted Payload
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.appender.AddNew(film.Name, film.Director); err != nil {
		switch {
		case errors.As(err, &driven.TypeFilmAlreadyExist):
			w.WriteHeader(http.StatusConflict)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	setHeaders(w)
	json.NewEncoder(w).Encode(appendResponse{
		Name:     film.Name,
		Director: film.Director,
		Release:  domain.Regular,
	})
}

func (s *server) addRegularFilm(w http.ResponseWriter, r *http.Request) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("request body read error : %w", err)
	}

	var film appendRequest
	if err := json.Unmarshal(reqBody, &film); err != nil || film.validate() {
		return NewClientError(err, http.StatusBadRequest, "Bad Request: Post payload cannot be deserialized")
	}

	if err := s.appender.AddRegular(film.Name, film.Director); err != nil {
		switch {
		case errors.As(err, &driven.TypeFilmAlreadyExist):
			return NewClientError(err, http.StatusConflict, "Status Conflict: Film Already Exist. Name must be unique!")
		default:
			return fmt.Errorf("unable to add film: %w", err)
		}
	}

	setHeaders(w)
	json.NewEncoder(w).Encode(appendResponse{
		Name:     film.Name,
		Director: film.Director,
		Release:  domain.Regular,
	})
	return nil
}

func (s *server) addOldFilm(w http.ResponseWriter, r *http.Request) error {
	return s.addFilm(w, r, s.appender.AddOld)
}

func (s *server) addFilm(w http.ResponseWriter, r *http.Request, fx func(name string, director string) error) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("request body read error : %w", err)
	}

	var film appendRequest
	if err := json.Unmarshal(reqBody, &film); err != nil {
		return NewClientError(err, http.StatusBadRequest, "Bad Request: Post payload cannot be deserialized")
	}

	if err := fx(film.Name, film.Director); err != nil {
		switch {
		case errors.As(err, &driven.TypeFilmAlreadyExist):
			return NewClientError(err, http.StatusConflict, "Status Conflict: Film Already Exist. Name must be unique!")
		default:
			return fmt.Errorf("unable to add film: %w", err)
		}
	}

	setHeaders(w)
	json.NewEncoder(w).Encode(appendResponse{
		Name:     film.Name,
		Director: film.Director,
		Release:  domain.Regular,
	})
	return nil
}
