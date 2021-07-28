package domain

import "fmt"

var ErrorUnknownReleaseType = fmt.Errorf("unknown release type must be one of the following releases, %v", releaseTypes)

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
	return f.Release.isValid()
}

func (r *release) isValid() error {
	switch *r {
	case Old, Regular, New:
		return nil
	}
	return ErrorUnknownReleaseType
}
