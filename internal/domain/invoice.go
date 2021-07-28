package domain

type (
	SEK  uint64
	Days uint16
)

const (
	PREMIUM = SEK(40)
	BASIC   = SEK(30)
)

type (
	Calculator func(days Days) SEK

	Rental struct {
		Film Film
		Days Days
	}

	RentalReturn struct {
		Rentals []Rental
	}

	RentalInvoice struct {
		RentalReturn
		Cost SEK
	}
)

func (req *RentalReturn) AddRental(film Film, days Days) {
	req.Rentals = append(req.Rentals, Rental{film, days})
}

var pricingStrategies = map[release]Calculator{
	New:     NewRelease,
	Regular: RegularRelease,
	Old:     OldRelease,
}

func getReleaseCalculator(release release) (Calculator, error) {
	if err := release.isValid(); err != nil {
		return nil, err
	}
	return pricingStrategies[release], nil
}

func (req RentalReturn) Invoice() (i RentalInvoice, e []error) {
	var cost = SEK(0)

	for _, r := range req.Rentals {
		if calc, err := getReleaseCalculator(r.Film.Release); err != nil {
			e = append(e, err)
		} else {
			cost += calc(r.Days)
		}
	}

	i = RentalInvoice{
		RentalReturn: req,
		Cost:         cost,
	}

	return i, e
}

func NewRelease(days Days) SEK {
	return SEK(uint64(days) * uint64(PREMIUM))
}

func RegularRelease(days Days) SEK {
	var gracePeriod = Days(3)
	return calculatePrice(days, gracePeriod)
}

func OldRelease(days Days) SEK {
	var gracePeriod = Days(5)
	return calculatePrice(days, gracePeriod)
}

func calculatePrice(days Days, gracePeriod Days) SEK {
	switch {
	case days == 0:
		return 0
	case days <= gracePeriod:
		return BASIC
	default:
		var excess = SEK(days.subtract(gracePeriod)) * BASIC
		return BASIC + excess
	}
}

func (d Days) subtract(deduct Days) Days {
	if deduct > d {
		return Days(0)
	}
	return d - deduct
}
