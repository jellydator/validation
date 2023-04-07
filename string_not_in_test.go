// Copyright 2016 Qiang Xue, Google LLC. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringNotIn(t *testing.T) {
	var tests = []struct {
		tag             string
		isCaseSensitive bool
		values          []string
		value           string
		err             string
	}{
		{"t0", true, []string{"A", "B"}, "", ""},
		{"t1", true, []string{"A", "B"}, "A", "must not be in list"},
		{"t2", true, []string{"A", "B"}, "B", "must not be in list"},
		{"t3", true, []string{"A", "B"}, "C", ""},
		{"t4", true, []string{}, "C", ""},
		{"t5", false, []string{"A", "B"}, "", ""},
		{"t6", false, []string{"A", "B"}, "A", "must not be in list"},
		{"t7", false, []string{"A", "B"}, "B", "must not be in list"},
		{"t8", false, []string{"A", "B"}, "C", ""},
		{"t9", false, []string{"A", "B"}, "a", "must not be in list"},
		{"t10", false, []string{"A", "B"}, "b", "must not be in list"},
		{"t11", false, []string{"A", "B"}, "c", ""},
		{"t12", false, []string{}, "c", ""},
	}

	for _, test := range tests {
		r := StringNotIn(test.isCaseSensitive, test.values...)
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
	}
}

func Test_StringNotInRule_Error(t *testing.T) {
	r := StringNotIn(true, "A", "B", "C")
	assert.Equal(t, "must not be in list", r.Validate("A").Error())
	r = r.Error("123")
	assert.Equal(t, "123", r.err.Message())
}

func TestStringNotInRule_ErrorObject(t *testing.T) {
	r := StringNotIn(true, "A", "B", "C")

	err := NewError("code", "abc")
	r = r.ErrorObject(err)

	assert.Equal(t, err, r.err)
	assert.Equal(t, err.Code(), r.err.Code())
	assert.Equal(t, err.Message(), r.err.Message())
}
