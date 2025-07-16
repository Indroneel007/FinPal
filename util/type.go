package util

const (
	rent     = "rent"
	food     = "food"
	travel   = "travel"
	salary   = "salary"
	savings  = "savings"
	bills    = "bills"
	medical  = "medical"
	shopping = "shopping"
	misc     = "misc"
)

func IsSupportedType(accountType string) bool {
	switch accountType {
	case rent, food, travel, salary, savings, bills, medical, shopping, misc:
		return true
	default:
		return false
	}
}
