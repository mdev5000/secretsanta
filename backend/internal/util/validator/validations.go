package validator

// @todo test package

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Assertion struct {
	ValidationName string
	Validation     interface{}
	Err            error
}

type Result struct {
	Results map[string]*ValueResult
}

func (r *Result) IsValid() bool {
	return len(r.Results) == 0
}

type ValueResult struct {
	Field  string
	Value  interface{}
	Errors []Assertion
}

type Results interface {
	AddError(field string, v interface{}, validationName string, err error)
}

func (r *Result) AddError(field string, v interface{}, validationName string, err error) {
	values, ok := r.Results[field]
	if !ok {
		values = &ValueResult{}
	}
	values.AddError(v, validationName, err)
	if !ok {
		r.Results[field] = values
	}
}

func (r *ValueResult) Error() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("validation failed for field '%s':\n", r.Field))
	for _, assertion := range r.Errors {
		b.WriteString("  ")
		b.WriteString(assertion.Err.Error())
		b.WriteString("\n")
	}
	return b.String()
}

func (r *ValueResult) AddError(v interface{}, validationName string, err error) {
	r.Errors = append(r.Errors, Assertion{
		ValidationName: validationName,
		Validation:     v,
		Err:            err,
	})
}

type Validator[T any] struct {
	ValidationName string
	Validate       func(Results, string, T)
}

func Regex(validationName string, err error, r *regexp.Regexp) Validator[string] {
	var v Validator[string]
	v = Validator[string]{
		ValidationName: validationName,
		Validate: func(results Results, field string, s string) {
			if !r.MatchString(s) {
				results.AddError(field, v, validationName, err)
			}
		},
	}
	return v
}

var (
	nameFieldRegex = regexp.MustCompile(`^[^{}$]+$`)

	NameFieldErr = errors.New("cannot contain {,},$")

	StringSearchField = Regex("nameField", NameFieldErr, nameFieldRegex)
)

func Validate[T any](results Results, field string, value T, validators ...Validator[T]) {
	for _, validator := range validators {
		validator.Validate(results, field, value)
	}
}

func ValidateSingle[T any](field string, value T, validators ...Validator[T]) error {
	results := &Result{Results: map[string]*ValueResult{}}
	Validate(results, field, value, validators...)
	vr, ok := results.Results[field]
	if !ok {
		return nil
	}
	return vr
}
