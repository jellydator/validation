// Copyright 2018 Qiang Xue, Google LLC. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import "strings"

// StringNotIn returns a validation rule that checks if a value is absent from the given list of values.
// An empty value is considered valid. Use the Required rule to make sure a value is not empty.
func StringNotIn(isCaseSensitive bool, values ...string) StringNotInRule {
	return StringNotInRule{
		isCaseSensitive: isCaseSensitive,
		elements:        values,
		err:             ErrNotInInvalid,
	}
}

// StringNotInRule is a validation rule that checks if a value is absent from the given list of values.
type StringNotInRule struct {
	isCaseSensitive bool
	elements        []string
	err             Error
}

// Validate checks if the given value is valid or not.
func (r StringNotInRule) Validate(value interface{}) error {
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
				return r.err
			}
		} else {
			if strings.EqualFold(e, valueAsString) {
				return r.err
			}
		}
	}

	return nil
}

// Error sets the error message for the rule.
func (r StringNotInRule) Error(message string) StringNotInRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r StringNotInRule) ErrorObject(err Error) StringNotInRule {
	r.err = err
	return r
}
