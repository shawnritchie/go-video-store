package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
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

	defer r.Body.Close()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//Handle Corrupted Payload
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var film appendRequest
	if err := json.Unmarshal(reqBody, &film); err != nil || film.validate() {
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appendResponse{
		Name:     film.Name,
		Director: film.Director,
		Release:  string(domain.New),
	})
}

func (s *server) addRegularFilm(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
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
		Release:  string(domain.Regular),
	})
	return nil
}

func (s *server) addFilm(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("request body read error : %w", err)
	}

	var film appendRequest
	if err := json.Unmarshal(reqBody, &film); err != nil || film.validate() {
		return NewClientError(err, http.StatusBadRequest, "Bad Request: Post payload cannot be deserialized")
	}

	pathParams := mux.Vars(r)
	strRelease, ok := pathParams["release"]
	if !ok {
		return NewClientError(err, http.StatusBadRequest, "Bad Request: missing release type within request. example: \"/catalogue/film/[new,regular,old]\"")
	}

	release, err := domain.ParseRelease(strRelease)
	if err != nil {
		return NewClientError(err, http.StatusBadRequest, "Bad Request: supported release types \"[new,regular,old]\"")
	}

	var fx func(name string, director string) error

	switch release {
	case domain.New:
		fx = s.appender.AddNew
	case domain.Regular:
		fx = s.appender.AddRegular
	case domain.Old:
		fx = s.appender.AddOld
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
		Release:  string(release),
	})
	return nil
}
