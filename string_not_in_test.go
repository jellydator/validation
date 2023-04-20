// Copyright 2016 Qiang Xue, Google LLC. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringNotIn(t *testing.T) {
	var v1 = "A"
	var v2 *string
	var v3 = "a"
	var tests = []struct {
		tag             string
		isCaseSensitive bool
		values          []string
		value           interface{}
		err             string
	}{
		{"t0", true, []string{"A", "B"}, "", ""},
		{"t1", true, []string{"A", "B"}, "A", "must not be in list"},
		{"t2", true, []string{"A", "B"}, "B", "must not be in list"},
		{"t3", true, []string{"A", "B"}, "C", ""},
		{"t4", true, []string{"A", "B"}, 4, "must be either a string or byte slice"},
		{"t5", true, []string{}, "C", ""},
		{"t6", true, []string{"A", "B"}, &v1, "must not be in list"},
		{"t7", true, []string{"A", "B"}, v2, ""},
		{"t8", true, []string{"A", "B"}, (*int)(nil), "must be either a string or byte slice"},
		{"t9", false, []string{"A", "B"}, "", ""},
		{"t10", false, []string{"A", "B"}, "A", "must not be in list"},
		{"t11", false, []string{"A", "B"}, "B", "must not be in list"},
		{"t12", false, []string{"A", "B"}, "C", ""},
		{"t13", false, []string{"A", "B"}, "a", "must not be in list"},
		{"t14", false, []string{"A", "B"}, "b", "must not be in list"},
		{"t15", false, []string{"A", "B"}, "c", ""},
		{"t16", false, []string{"A", "B"}, 4, "must be either a string or byte slice"},
		{"t17", false, []string{}, "c", ""},
		{"t18", false, []string{"A", "B"}, &v3, "must not be in list"},
		{"t19", false, []string{"A", "B"}, v2, ""},
	}

	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			r := StringNotIn(test.isCaseSensitive, test.values...)
			err := r.Validate(test.value)
			assertError(t, test.err, err, test.tag)
		})
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
