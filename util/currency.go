package util

const (
	USD    = "USD"
	Euros  = "Euros"
	Rupees = "Rupees"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, Euros, Rupees:
		return true
	default:
		return false
	}
}
