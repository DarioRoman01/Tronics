package handlers

import "gopkg.in/go-playground/validator.v9"

var v = validator.New()

type ProductValidator struct {
	validator *validator.Validate
}

// validate request payload for products
func (p *ProductValidator) Validate(i interface{}) error {
	return p.validator.Struct(i)
}

type UserValidator struct {
	validator *validator.Validate
}

// validate request payload for users
func (u *UserValidator) Validate(i interface{}) error {
	return u.validator.Struct(i)
}
