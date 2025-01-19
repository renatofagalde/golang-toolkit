package toolkit

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	BRL = "BRL"
)

func IsSuppertedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, BRL:
		return true
	}
	return false
}
