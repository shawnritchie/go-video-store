package domain

type (
	release string

	Film struct {
		Name     string
		Director string
		Release  release
	}
)

const (
	New     release = "New"
	Regular         = "Regular"
	Old             = "Old"
)

var releaseTypes = []release{New, Regular, Old}

func (f *Film) IsValid() error {
	var errors []error
	if f.Name == "" {
		errors = append(errors, EmptyFilmNameError)
	}

	if f.Director == "" {
		errors = append(errors, EmptyFilmDirectorError)
	}

	if err := f.Release.isValid(); err != nil {
		errors = append(errors, EmptyFilmDirectorError)
	}

	if len(errors) == 0 {
		return nil
	} else {
		var ret InvalidFilmError = errors
		return &ret
	}
}

func (r *release) isValid() error {
	switch *r {
	case Old, Regular, New:
		return nil
	}
	return UnknownReleaseError
}
