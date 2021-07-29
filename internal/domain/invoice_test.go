package domain

import (
	"testing"
)

var newFilm = Film{"Loki", "Marvel", New}
var regularFilm = Film{"Loki", "Marvel", Regular}
var oldFilm = Film{"Loki", "Marvel", Old}

func TestNewFilmPricing(t *testing.T) {
	tests := []struct {
		film          Film
		days          Days
		expectedPrice SEK
	}{
		{newFilm, 0, 0},
		{newFilm, 1, PREMIUM},
		{newFilm, 2, PREMIUM * 2},
		{newFilm, 4, PREMIUM * 4},
		{newFilm, 10, PREMIUM * 10},
		{regularFilm, 0, 0},
		{regularFilm, 1, BASIC},
		{regularFilm, 3, BASIC},
		{regularFilm, 4, BASIC * 2},
		{regularFilm, 10, BASIC * 8},
		{oldFilm, 0, 0},
		{oldFilm, 1, BASIC},
		{oldFilm, 5, BASIC},
		{oldFilm, 6, BASIC * 2},
		{oldFilm, 10, BASIC * 6},
	}
	for _, test := range tests {
		t.Run("Pricing Test New release", func(t *testing.T) {
			if calc, err := getReleaseCalculator(test.film.Release); err != nil {
				t.Error(err)
			} else {
				var cost = calc(test.days)
				if cost != test.expectedPrice {
					t.Errorf("calculated cost of %d didn't match expect price %d", cost, test.expectedPrice)
				}
			}

		})
	}
}

func TestInvoicing(t *testing.T) {
	var duration = Days(5)

	var request = RentalReturn{
		Rentals: []Rental{
			{newFilm, duration},
			{regularFilm, duration},
			{oldFilm, duration},
		},
	}

	if invoice, err := request.Invoice(); err != nil {
		t.Error(err)
	} else {
		var expectedPrice = PREMIUM*SEK(5) + BASIC*SEK(3) + BASIC
		if invoice.Cost != expectedPrice {
			t.Errorf("calculated cost of %d didn't match expect price %d", invoice.Cost, expectedPrice)
		}
	}
}

func TestCorruptedRentalRequest(t *testing.T) {
	var corruptedFilm = Film{"Boki", "DC", release("Corrupted")}
	var duration = Days(5)

	var request = RentalReturn{
		Rentals: []Rental{
			{corruptedFilm, duration},
			{corruptedFilm, duration},
			{corruptedFilm, duration},
		},
	}

	if _, errors := request.Invoice(); errors != nil && len(errors) != 3 {
		t.Errorf("was expeecting 3 errors for 3 none existant release types")
	}
}
