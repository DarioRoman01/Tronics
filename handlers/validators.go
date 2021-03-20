package handlers

import "gopkg.in/go-playground/validator.v9"

var v = validator.New()

type ProductValidator struct {
	validator *validator.Validate
}

func (p *ProductValidator) Validate(i interface{}) error {
	return p.validator.Struct(i)
}
