package service

import (
	"errors"
	"github.com/shawnritchie/go-video-store/internal/domain"
	"github.com/shawnritchie/go-video-store/internal/port/driven"
	"github.com/shawnritchie/go-video-store/internal/port/driver"
)

type (
	StoreService struct {
		finder   driver.Queryable
		appender driver.Insertable
	}
)

func New(finder driver.Queryable, appender driver.Insertable) *StoreService {
	return &StoreService{
		finder,
		appender,
	}
}

func (svc *StoreService) Find(name string) (*domain.Film, error) {
	return svc.finder.FindBy(name)
}

func (svc *StoreService) AddNew(name string, director string) error {
	return svc.addFilm(domain.Film{Name: name, Director: director, Release: domain.New})
}

func (svc *StoreService) AddRegular(name string, director string) error {
	return svc.addFilm(domain.Film{Name: name, Director: director, Release: domain.Regular})
}

func (svc *StoreService) AddOld(name string, director string) error {
	return svc.addFilm(domain.Film{Name: name, Director: director, Release: domain.Old})
}

func (svc *StoreService) Invoice(request []driven.FilmReturn) (*domain.RentalInvoice, error) {
	rentalRequest, invalidReq := svc.validateFilmReturn(request)
	if len(invalidReq) > 0 {
		return nil, &invalidReq
	}

	if invoice, errors := rentalRequest.Invoice(); errors != nil {
		error := driven.InvalidRentalRequestError(errors)
		return nil, &error
	} else {
		return &invoice, nil
	}
}

func (svc *StoreService) addFilm(film domain.Film) error {
	if err := film.IsValid(); err != nil {
		return err
	}

	if _, err := svc.finder.FindBy(film.Name); err != nil {
		if errors.As(err, &driven.TypeFilmNotFound) {
			return svc.appender.Insert(film)
		}
		return err
	}

	return nil
}

func (svc *StoreService) validateFilmReturn(request []driven.FilmReturn) (req domain.RentalReturn, invalidReq driven.InvalidRentalRequestError) {
	invalidReq = driven.InvalidRentalRequestError{}
	for _, rental := range request {
		if film, err := svc.finder.FindBy(rental.FilmName); err != nil {
			invalidReq.Append(err)
		} else {
			req.AddRental(*film, domain.Days(rental.Days))
		}
	}
	return req, invalidReq
}
