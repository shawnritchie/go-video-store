package driven

import "fmt"

type (
	InvalidRentalRequestError []error

	FilmNotFoundError struct {
		Name string
	}

	FilmAlreadyExistError struct {
		Name string
	}
)

var (
	TypeFilmNotFound         *FilmNotFoundError
	TypeFilmAlreadyExist     *FilmAlreadyExistError
	TypeInvalidRentalRequest *InvalidRentalRequestError
)

func (e *FilmNotFoundError) Error() string {
	return fmt.Sprintf("film: %q was not found", e.Name)
}

func (e *FilmAlreadyExistError) Error() string {
	return fmt.Sprintf("film: %q already exists", e.Name)
}

func (e *InvalidRentalRequestError) Error() (errMsg string) {
	errMsg = fmt.Sprintf("%d errors encountered\n", len(*e))
	for _, err := range *e {
		errMsg += fmt.Sprintf("- %s\n", err.Error())
	}
	return errMsg
}

func (e *InvalidRentalRequestError) Append(err error) {
	*e = append(*e, err)
}
