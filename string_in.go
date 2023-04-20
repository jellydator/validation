// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"strings"
)

// StringIn returns a validation rule that checks if a value can be found in the given list of values.
// == or strings.EqualFold will be used to determine if two values are equal.
// An empty value is considered valid. Use the Required rule to make sure a value is not empty.
func StringIn(isCaseSensitive bool, values ...string) StringInRule {
	return StringInRule{
		isCaseSensitive: isCaseSensitive,
		elements:        values,
		err:             ErrInInvalid,
	}
}

// StringInRule is a validation rule that validates if a value can be found in the given list of values.
type StringInRule struct {
	isCaseSensitive bool
	elements        []string
	err             Error
}

// Validate checks if the given value is valid or not.
func (r StringInRule) Validate(value interface{}) error {
	_, isStringPtr := value.(*string)
	indirectValue, isNil := Indirect(value)

	if isNil && isStringPtr {
		return nil
	}
	valueAsString, err := EnsureString(indirectValue)
	if err != nil {
		return err
	}
	if IsEmpty(indirectValue) {
		return nil
	}

	for _, e := range r.elements {
		if r.isCaseSensitive {
			if e == valueAsString {
				return nil
			}
		} else {
			if strings.EqualFold(e, valueAsString) {
				return nil
			}
		}
	}

	return r.err
}

// Error sets the error message for the rule.
func (r StringInRule) Error(message string) StringInRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r StringInRule) ErrorObject(err Error) StringInRule {
	r.err = err
	return r
}
