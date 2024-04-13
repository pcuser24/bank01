package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/user2410/simplebank/util"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		for _, c := range util.Currencies {
			if c == currency {
				return true
			}
		}
	}
	return false
}
