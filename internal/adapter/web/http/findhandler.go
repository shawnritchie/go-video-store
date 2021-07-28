package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shawnritchie/go-video-store/internal/port/driven"
	"net/http"
)

type (
	findResponse struct {
		Name     string `json:"name"`
		Director string `json:"director"`
		Release  string `json:"release"`
	}
)

func (s *server) findFilm(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	filmName := query.Get("name")

	if filmName == "" {
		return NewClientError(nil, http.StatusBadRequest, "Bad Request: Expected query parameter \"name\" in url")
	}

	film, err := s.finder.Find(filmName)
	if errors.As(err, &driven.TypeFilmNotFound) {
		return NewClientError(nil, http.StatusNotFound, fmt.Sprintf("Film Not Found: Film %q not found", filmName))
	} else if err != nil {
		return err
	}

	setHeaders(w)
	json.NewEncoder(w).Encode(findResponse{
		Name:     film.Name,
		Director: film.Director,
		Release:  string(film.Release),
	})
	return nil
}
