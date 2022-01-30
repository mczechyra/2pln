package api

type Currency int

const (
	EUR Currency = iota
	USD
	GBP
	CHF
)

func (c Currency) String() string {
	switch c {
	case EUR:
		return "EUR"
	case USD:
		return "USD"
	case GBP:
		return "GBP"
	case CHF:
		return "CHF"
	}
	return "EUR"
}
