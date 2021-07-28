package inmem

import (
	"errors"
	"github.com/shawnritchie/go-video-store/internal/domain"
	"github.com/shawnritchie/go-video-store/internal/port/driven"
	"github.com/shawnritchie/go-video-store/internal/port/driver"
	"testing"
)

//Array Declaration
var catalogue = StoreCatalogue{
	domain.Film{Name: "Matrix 11", Director: "Dwight", Release: domain.New},
	domain.Film{Name: "Spider Man", Director: "Dwight", Release: domain.Regular},
	domain.Film{Name: "Spider Man 2", Director: "Dwight", Release: domain.Regular},
	domain.Film{Name: "Out of Africa", Director: "Dwight", Release: domain.Old},
}

var repo driver.Catalogue = &catalogue

func TestFindFilm(t *testing.T) {
	var find = catalogue[0]
	if found, err := repo.FindBy(find.Name); err != nil {
		t.Error(err)
	} else if find != *found {
		t.Errorf("searched for %q but got %q", find.Name, found.Name)
	}
}

func TestFindFilm_FilmNotFoundError(t *testing.T) {
	found, err := repo.FindBy("Black Widow")

	if err == nil || found != nil {
		t.Errorf("was expecting film to be nil and err to be FilmNotFoundError")
	}

	if !errors.As(err, &driven.TypeFilmNotFound) {
		t.Errorf("was expecting TypeFilmNotFound error but got %#v", err)
	}
}

func TestAddFilm(t *testing.T) {
	var newFilm = domain.Film{
		Name:     "Loki",
		Director: "Marvel",
		Release:  domain.New,
	}

	if err := repo.Insert(newFilm); err != nil {
		t.Errorf("was expecting film to be inserted succesfully film: %#v but failed with %v", newFilm, err)
	}

	if _, err := repo.FindBy(newFilm.Name); err != nil {
		t.Error(err)
	}
}
