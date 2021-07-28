package service

import (
	"github.com/shawnritchie/go-video-store/internal/adapter/repository/inmem"
	"github.com/shawnritchie/go-video-store/internal/domain"
	"github.com/shawnritchie/go-video-store/internal/port/driven"
	"github.com/shawnritchie/go-video-store/internal/port/driver"
	"testing"
)

var films = []domain.Film{
	{"Matrix 11", "Dwight", domain.New},
	{"Spider Man", "Dwight", domain.Regular},
	{"Spider Man 2", "Dwight", domain.Regular},
	{"Out of Africa", "Dwight", domain.Old},
}

func setupCatalogue() (*inmem.StoreCatalogue, driver.Catalogue) {
	cat := inmem.StoreCatalogue(films)
	return &cat, &cat
}

func TestStoreService_AddNewFilm(t *testing.T) {
	repo, catalogue := setupCatalogue()
	service := New(catalogue, catalogue)

	newFilm := domain.Film{Name: "Loki", Director: "Marvel", Release: domain.New}
	service.AddNew(newFilm.Name, newFilm.Director)

	if (*repo)[4] != newFilm {
		t.Errorf("film %+v hasn't been added to catalogue", newFilm)
	}
}

func TestStoreService_AddRegularFilm(t *testing.T) {
	repo, catalogue := setupCatalogue()
	service := New(catalogue, catalogue)

	regularFilm := domain.Film{Name: "Loki", Director: "Marvel", Release: domain.Regular}
	service.AddRegular(regularFilm.Name, regularFilm.Director)

	if (*repo)[4] != regularFilm {
		t.Errorf("film %+v hasn't been added to catalogue", regularFilm)
	}
}

func TestStoreService_AddOldFilm(t *testing.T) {
	repo, catalogue := setupCatalogue()
	service := New(catalogue, catalogue)

	oldFilm := domain.Film{Name: "Loki", Director: "Marvel", Release: domain.Old}
	service.AddOld(oldFilm.Name, oldFilm.Director)

	if (*repo)[4] != oldFilm {
		t.Errorf("film %+v hasn't been added to catalogue", oldFilm)
	}
}

func TestStoreService_FindByName(t *testing.T) {
	repo, catalogue := setupCatalogue()
	service := New(catalogue, catalogue)

	lookFor := (*repo)[0]
	if found, err := service.Find(lookFor.Name); err != nil {
		t.Error(err)
	} else if *found != lookFor {
		t.Errorf("found wrong film %#v was looking for %#v", found, lookFor)
	}
}

func TestStoreService_StoreReturn(t *testing.T) {
	_, catalogue := setupCatalogue()
	service := New(catalogue, catalogue)

	duration := uint16(5)
	if invoice, err := service.Invoice(mapFilmReturn(films, duration)); err != nil {
		t.Error(err)
	} else {
		for i, rental := range invoice.Rentals {
			if rental.Film != films[i] {
				t.Errorf("film %#v is missing from invoice", films[i])
			} else if rental.Days != domain.Days(duration) {
				t.Errorf("film %#v has been incorrectly invoiced billed duration %d actual duration", rental.Days, duration)
			}
		}

		if invoice.Cost <= domain.SEK(0) {
			t.Errorf("incorrectly invoiced amount")
		}
	}
}

func mapFilmReturn(films []domain.Film, duration uint16) (ret []driven.FilmReturn) {
	for _, f := range films {
		ret = append(ret, driven.FilmReturn{FilmName: f.Name, Days: duration})
	}
	return ret
}
