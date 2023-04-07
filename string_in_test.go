// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringIn(t *testing.T) {
	tests := []struct {
		tag             string
		isCaseSensitive bool
		values          []string
		value           string
		err             string
	}{
		{"t0", true, []string{"A", "B"}, "", ""},
		{"t1", true, []string{"A", "B"}, "A", ""},
		{"t2", true, []string{"A", "B"}, "B", ""},
		{"t3", true, []string{"A", "B"}, "C", "must be a valid value"},
		{"t4", true, []string{}, "C", "must be a valid value"},
		{"t5", false, []string{"A", "B"}, "", ""},
		{"t6", false, []string{"A", "B"}, "A", ""},
		{"t7", false, []string{"A", "B"}, "B", ""},
		{"t8", false, []string{"A", "B"}, "C", "must be a valid value"},
		{"t9", false, []string{"A", "B"}, "a", ""},
		{"t10", false, []string{"A", "B"}, "b", ""},
		{"t11", false, []string{"A", "B"}, "c", "must be a valid value"},
		{"t12", false, []string{}, "c", "must be a valid value"},
	}

	for _, test := range tests {
		r := StringIn(test.isCaseSensitive, test.values...)
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
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
