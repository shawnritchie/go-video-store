package service

import (
	"github.com/shawnritchie/go-video-store/internal/adapter/repository/inmem"
	"github.com/shawnritchie/go-video-store/internal/domain"
	"github.com/shawnritchie/go-video-store/internal/port/driven"
	"github.com/shawnritchie/go-video-store/internal/port/driver"
	"testing"
)

var films = []domain.Film{
	{Name: "Matrix 11", Director: "Dwight", Release: domain.New},
	{Name: "Spider Man", Director: "Dwight", Release: domain.Regular},
	{Name: "Spider Man 2", Director: "Dwight", Release: domain.Regular},
	{Name: "Out of Africa", Director: "Dwight", Release: domain.Old},
}

type spyCatalogue struct {
	findBy func(name string) (*domain.Film, error)
	insert func(film domain.Film) error
}

func (s *spyCatalogue) FindBy(name string) (*domain.Film, error) {
	return s.findBy(name)
}

func (s *spyCatalogue) Insert(film domain.Film) error {
	return s.insert(film)
}

func newSpyCatalogue(
	findBy func(name string) (*domain.Film, error),
	insert func(film domain.Film) error) *spyCatalogue {
	return &spyCatalogue{
		findBy: findBy,
		insert: insert,
	}
}

func setupCatalogue() driver.Catalogue {
	cat := inmem.StoreCatalogue(films)
	return &cat
}

func mockFindByError(err error) func(name string) (*domain.Film, error) {
	return func(name string) (*domain.Film, error) {
		return nil, err
	}
}

func TestStoreService_AddNewFilm(t *testing.T) {
	hasBeenInvoked := false
	newFilm := domain.Film{Name: "Loki", Director: "Marvel", Release: domain.New}

	catalogue := newSpyCatalogue(
		mockFindByError(&driven.FilmNotFoundError{Name: newFilm.Name}),
		func(film domain.Film) error {
			hasBeenInvoked = true
			if film != newFilm {
				t.Errorf("film %+v hasn't been added to catalogue", newFilm)
			}
			return nil
		})

	service := New(catalogue, catalogue)
	service.AddNew(newFilm.Name, newFilm.Director)

	if !hasBeenInvoked {
		t.Errorf("film %+v hasn't been added to catalogue", newFilm)
	}
}

func TestAddFilm(t *testing.T) {
	tests := []struct {
		testName     string
		insertedFilm domain.Film
		addFx        func(s *StoreService, f domain.Film)
	}{
		{
			testName:     "addNewFilmTest",
			insertedFilm: domain.Film{Name: "Loki", Director: "Marvel", Release: domain.New},
			addFx:        func(s *StoreService, f domain.Film) { s.AddNew(f.Name, f.Director) },
		},
		{
			testName:     "addRegularFilmTest",
			insertedFilm: domain.Film{Name: "Loki", Director: "Marvel", Release: domain.Regular},
			addFx:        func(s *StoreService, f domain.Film) { s.AddRegular(f.Name, f.Director) },
		},
		{
			testName:     "addOldFilmTest",
			insertedFilm: domain.Film{Name: "Loki", Director: "Marvel", Release: domain.Old},
			addFx:        func(s *StoreService, f domain.Film) { s.AddOld(f.Name, f.Director) },
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			hasBeenInvoked := false

			catalogue := newSpyCatalogue(
				mockFindByError(&driven.FilmNotFoundError{Name: test.insertedFilm.Name}),
				func(film domain.Film) error {
					hasBeenInvoked = true
					if film != test.insertedFilm {
						t.Errorf("film %#v hasn't been added to catalogue", test.insertedFilm)
					}
					return nil
				})

			service := New(catalogue, catalogue)
			test.addFx(service, test.insertedFilm)

			if !hasBeenInvoked {
				t.Errorf("film %#v hasn't been added to catalogue", test.insertedFilm)
			}
		})
	}
}

func TestStoreService_FindByName(t *testing.T) {
	searchFor := domain.Film{Name: "Loki", Director: "Marvel", Release: domain.New}
	hasBeenInvoked := false

	catalogue := newSpyCatalogue(
		func(name string) (*domain.Film, error) {
			hasBeenInvoked = true
			if name != searchFor.Name {
				t.Errorf("looking for wrong film expected search was %q, but search for %q", searchFor, name)
			}
			return &searchFor, nil
		},
		nil)

	service := New(catalogue, catalogue)
	service.Find(searchFor.Name)

	if !hasBeenInvoked {
		t.Errorf("findBy hasn't been invoked")
	}
}

func TestStoreService_StoreReturn(t *testing.T) {
	catalogue := setupCatalogue()
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
