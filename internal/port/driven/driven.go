package driven

import (
	"github.com/shawnritchie/go-video-store/internal/domain"
)

type (
	FilmReturn struct {
		FilmName string
		Days     uint16
	}
)

type (
	FilmFinder interface {
		Find(name string) (*domain.Film, error)
	}

	FilmAppender interface {
		AddNew(name string, director string) error
		AddRegular(name string, director string) error
		AddOld(name string, director string) error
	}

	FilmInvoicer interface {
		Invoice(request []FilmReturn) (*domain.RentalInvoice, error)
	}
)
