package domain

import (
	"testing"
)

func TestValidReleases(t *testing.T) {
	releases := []release{
		New,
		Regular,
		Old,
	}
	for _, release := range releases {
		t.Run(string(release), func(t *testing.T) {
			if err := release.isValid(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestInvalidRelease(t *testing.T) {
	var unknownRelease = release("Disney")
	if err := unknownRelease.isValid(); err == nil {
		t.Errorf("unknown release %q should return an error", unknownRelease)
	}
}
