package driven

import (
	"fmt"
	"github.com/shawnritchie/go-video-store/internal/domain"
)

type (
	FilmReturn struct {
		FilmName string
		Days     uint16
	}

	InvalidRentalRequest []error
)

type (
	FilmNotFoundError struct {
		Name string
	}

	FilmAlreadyExist struct {
		Name string
	}

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

var (
	TypeFilmNotFound         *FilmNotFoundError
	TypeFilmAlreadyExist     *FilmAlreadyExist
	TypeInvalidRentalRequest *InvalidRentalRequest
)

func (e *FilmNotFoundError) Error() string {
	return fmt.Sprintf("film: %q was not found", e.Name)
}

func (e *FilmAlreadyExist) Error() string {
	return fmt.Sprintf("film: %q already exists", e.Name)
}

func (e *InvalidRentalRequest) Error() (errMsg string) {
	errMsg = fmt.Sprintf("%d errors encountered\n", len(*e))
	for _, err := range *e {
		errMsg += fmt.Sprintf("- %s\n", err.Error())
	}
	return errMsg
}

func (e *InvalidRentalRequest) Append(err error) {
	*e = append(*e, err)
}
