package validation

import (
	"fmt"
	"strings"
)

type ValidateError []error

func (e ValidateError) Error() string {
	var res strings.Builder
	for _, err := range e {
		res.WriteString(fmt.Sprintf("%s\n", err.Error()))
	}
	return res.String()
}

func (e ValidateError) ToArrayErrors() []string {
	var res []string
	for _, err := range e {
		res = append(res, err.Error())
	}
	return res
}

type Validator[T any] struct {
	errors ValidateError
	object T
}

func NewValidator[T any](object T) *Validator[T] {
	return &Validator[T]{
		object: object,
	}
}

func (v *Validator[T]) Between(getField func(T) int, start, end int, errMsg string) *Validator[T] {
	return v.Must(func(t T) bool {
		return getField(t) >= start && getField(t) <= end
	}, errMsg)
}

func (v *Validator[T]) Must(f func(T) bool, errMsg string) *Validator[T] {
	if !f(v.object) {
		v.errors = append(v.errors, fmt.Errorf(errMsg))
	}
	return v
}

func (v *Validator[T]) Validate() error {
	if len(v.errors) == 0 {
		return nil
	}
	return &v.errors
}
