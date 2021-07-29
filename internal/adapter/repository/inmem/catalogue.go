package inmem

import (
	"github.com/shawnritchie/go-video-store/internal/domain"
	"github.com/shawnritchie/go-video-store/internal/port/driven"
)

type (
	StoreCatalogue []domain.Film
)

func (cat *StoreCatalogue) FindBy(name string) (*domain.Film, error) {
	for _, film := range *cat {
		if film.Name == name {
			return &film, nil
		}
	}

	return nil, &driven.FilmNotFoundError{Name: name}
}

func (cat *StoreCatalogue) Insert(film domain.Film) error {
	/*
	 *Adapaters should be dumb variance rules should be part of the service
	 */
	*cat = append(*cat, film)
	return nil
}
