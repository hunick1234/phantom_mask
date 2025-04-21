package utils

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("timeformat", validateTimeFormat)
		v.RegisterValidation("dayofweek", validateDayOfWeek)
		v.RegisterValidation("sortby", validateSortBy)
		v.RegisterValidation("comparison", validateComparison)
	}
}

func validateTimeFormat(fl validator.FieldLevel) bool {
	_, err := time.Parse("15:04", fl.Field().String())
	return err == nil
}

func validateDayOfWeek(fl validator.FieldLevel) bool {
	day := fl.Field().String()
	validDays := map[string]bool{
		"Mon": true, "Tue": true, "Wed": true, "Thur": true, "Fri": true, "Sat": true, "Sun": true,
	}
	return validDays[day]
}

func validateSortBy(fl validator.FieldLevel) bool {
	sort := strings.ToLower(fl.Field().String())
	return sort == "price" || sort == "name"
}


func validateComparison(fl validator.FieldLevel) bool {
	value := strings.ToLower(fl.Field().String())
	return value == "more" || value == "less"
}