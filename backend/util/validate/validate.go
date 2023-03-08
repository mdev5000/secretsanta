package validate

import "github.com/go-playground/validator/v10"

type Validate = validator.Validate

var validation *Validate

func init() {
	validation = validator.New()
}

func Get() *Validate {
	return validation
}
