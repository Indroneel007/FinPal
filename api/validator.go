package api

import (
	"examples/SimpleBankProject/util"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	currency, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return util.IsSupportedCurrency(currency)
}
