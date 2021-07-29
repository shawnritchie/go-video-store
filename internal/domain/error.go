package domain

import "fmt"

type (
	InvalidFilmError []error
)

var (
	UnknownReleaseError    = fmt.Errorf("unknown release type must be one of the following releases, %v", releaseTypes)
	EmptyFilmNameError     = fmt.Errorf("film name cannot be empty")
	EmptyFilmDirectorError = fmt.Errorf("film director cannot be empty")

	TypeInvalidFilm *InvalidFilmError
)

func (e *InvalidFilmError) Error() (errMsg string) {
	errMsg = fmt.Sprintf("%d errors encountered\n", len(*e))
	for _, err := range *e {
		errMsg += fmt.Sprintf("- %s\n", err.Error())
	}
	return errMsg
}

func (e *InvalidFilmError) Append(err error) {
	*e = append(*e, err)
}
