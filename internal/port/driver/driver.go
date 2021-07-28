package driver

import (
	"github.com/shawnritchie/go-video-store/internal/domain"
)

type (
	Queryable interface {
		FindBy(name string) (*domain.Film, error)
	}

	Insertable interface {
		Insert(film domain.Film) error
	}

	Catalogue interface {
		Queryable
		Insertable
	}
)
