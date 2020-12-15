package handler

import (
	"github.com/go-playground/validator/v10"
	"time"
)

func CustomFuncInList(fl validator.FieldLevel) bool {
	ls := []string{
		"list1",
		"list2",
		"list3",
		"list4",
	}
	for _, v := range ls {
		if v == fl.Field().String() {
			return true
		}
	}
	return false
}

func DateFormat(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}
	return true
}
