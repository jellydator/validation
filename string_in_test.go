// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringIn(t *testing.T) {
	var v1 = "A"
	var v2 *string
	var v3 = "a"
	tests := []struct {
		tag             string
		isCaseSensitive bool
		values          []string
		value           interface{}
		err             string
	}{
		{"t0", true, []string{"A", "B"}, "", ""},
		{"t1", true, []string{"A", "B"}, "A", ""},
		{"t2", true, []string{"A", "B"}, "B", ""},
		{"t3", true, []string{"A", "B"}, "C", "must be a valid value"},
		{"t4", true, []string{"A", "B"}, 4, "must be either a string or byte slice"},
		{"t5", true, []string{}, "C", "must be a valid value"},
		{"t6", true, []string{"A", "B"}, &v1, ""},
		{"t7", true, []string{"A", "B"}, v2, ""},
		{"t8", true, []string{"A", "B"}, (*int)(nil), "must be either a string or byte slice"},
		{"t9", false, []string{"A", "B"}, "", ""},
		{"t10", false, []string{"A", "B"}, "A", ""},
		{"t11", false, []string{"A", "B"}, "B", ""},
		{"t12", false, []string{"A", "B"}, "C", "must be a valid value"},
		{"t13", false, []string{"A", "B"}, "a", ""},
		{"t14", false, []string{"A", "B"}, "b", ""},
		{"t15", false, []string{"A", "B"}, "c", "must be a valid value"},
		{"t16", false, []string{"A", "B"}, 4, "must be either a string or byte slice"},
		{"t17", false, []string{}, "c", "must be a valid value"},
		{"t18", false, []string{"A", "B"}, &v3, ""},
		{"t19", false, []string{"A", "B"}, v2, ""},
	}

	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			r := StringIn(test.isCaseSensitive, test.values...)
			err := r.Validate(test.value)
			assertError(t, test.err, err, test.tag)
		})
	}
}

func Test_StringInRule_Error(t *testing.T) {
	r := StringIn(true, "A", "B", "C")
	val := "D"
	assert.Equal(t, "must be a valid value", r.Validate(val).Error())
	r = r.Error("123")
	assert.Equal(t, "123", r.err.Message())
}

func TestStringInRule_ErrorObject(t *testing.T) {
	r := StringIn(true, "A", "B", "C")

	err := NewError("code", "abc")
	r = r.ErrorObject(err)

	assert.Equal(t, err, r.err)
	assert.Equal(t, err.Code(), r.err.Code())
	assert.Equal(t, err.Message(), r.err.Message())
}
