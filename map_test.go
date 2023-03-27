package validation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	var m0 map[string]interface{}
	m1 := map[string]interface{}{"A": "abc", "B": "xyz", "c": "abc", "D": (*string)(nil), "F": (*String123)(nil), "H": []string{"abc", "abc"}, "I": map[string]string{"foo": "abc"}}
	m2 := map[string]interface{}{"E": String123("xyz"), "F": (*String123)(nil)}
	m3 := map[string]interface{}{"M3": Model3{}}
	m4 := map[string]interface{}{"M3": Model3{A: "abc"}}
	m5 := map[string]interface{}{"A": "internal", "B": ""}
	m6 := map[int]string{11: "abc", 22: "xyz"}
	m7 := map[string]string{"xyz": "abc", "abc": "xyz"}
	m8 := map[string]string{"11": "abc", "222": "xyz"}
	tests := []struct {
		tag          string
		model        interface{}
		keyRules     []Rule
		valueRules   []Rule
		distinctKeys []*KeyRules
		err          string
	}{
		// empty rules
		{"t1.1", m1, nil, nil, []*KeyRules{}, ""},
		{"t1.2", m1, nil, nil, []*KeyRules{Key("A"), Key("B")}, ""},
		// normal rules
		{"t2.1", m1, nil, nil, []*KeyRules{Key("A", &validateAbc{}), Key("B", &validateXyz{})}, ""},
		{"t2.2", m1, nil, nil, []*KeyRules{Key("A", &validateXyz{}), Key("B", &validateAbc{})}, "A: error xyz; B: error abc."},
		{"t2.3", m1, nil, nil, []*KeyRules{Key("A", &validateXyz{}), Key("c", &validateXyz{})}, "A: error xyz; c: error xyz."},
		{"t2.4", m1, nil, nil, []*KeyRules{Key("D", Length(0, 5))}, ""},
		{"t2.5", m1, nil, nil, []*KeyRules{Key("F", Length(0, 5))}, ""},
		{"t2.6", m1, nil, nil, []*KeyRules{Key("H", Each(&validateAbc{})), Key("I", Each(&validateAbc{}))}, ""},
		{"t2.7", m1, nil, nil, []*KeyRules{Key("H", Each(&validateXyz{})), Key("I", Each(&validateXyz{}))}, "H: (0: error xyz; 1: error xyz.); I: (foo: error xyz.)."},
		{"t2.8", m1, nil, nil, []*KeyRules{Key("I", Map(Key("foo", &validateAbc{})))}, ""},
		{"t2.9", m1, nil, nil, []*KeyRules{Key("I", Map(Key("foo", &validateXyz{})))}, "I: (foo: error xyz.)."},
		// non-map value
		{"t3.1", &m1, nil, nil, []*KeyRules{}, ""},
		{"t3.2", nil, nil, nil, []*KeyRules{}, ErrNotMap.Error()},
		{"t3.3", m0, nil, nil, []*KeyRules{}, ""},
		{"t3.4", &m0, nil, nil, []*KeyRules{}, ""},
		{"t3.5", 123, nil, nil, []*KeyRules{}, ErrNotMap.Error()},
		// invalid key spec
		{"t4.1", m1, nil, nil, []*KeyRules{Key(123)}, "123: key not the correct type."},
		{"t4.2", m1, nil, nil, []*KeyRules{Key("X")}, "X: required key is missing."},
		{"t4.3", m1, nil, nil, []*KeyRules{Key("X").Optional()}, ""},
		// non-string keys
		{"t5.1", m6, nil, nil, []*KeyRules{Key(11, &validateAbc{}), Key(22, &validateXyz{})}, ""},
		{"t5.2", m6, nil, nil, []*KeyRules{Key(11, &validateXyz{}), Key(22, &validateAbc{})}, "11: error xyz; 22: error abc."},
		// validatable value
		{"t6.1", m2, nil, nil, []*KeyRules{Key("E")}, "E: error 123."},
		{"t6.2", m2, nil, nil, []*KeyRules{Key("E", Skip)}, ""},
		{"t6.3", m2, nil, nil, []*KeyRules{Key("E", Skip.When(true))}, ""},
		{"t6.4", m2, nil, nil, []*KeyRules{Key("E", Skip.When(false))}, "E: error 123."},
		// Required, NotNil
		{"t7.1", m2, nil, nil, []*KeyRules{Key("F", Required)}, "F: cannot be blank."},
		{"t7.2", m2, nil, nil, []*KeyRules{Key("F", NotNil)}, "F: is required."},
		{"t7.3", m2, nil, nil, []*KeyRules{Key("F", Skip, Required)}, ""},
		{"t7.4", m2, nil, nil, []*KeyRules{Key("F", Skip, NotNil)}, ""},
		{"t7.5", m2, nil, nil, []*KeyRules{Key("F", Skip.When(true), Required)}, ""},
		{"t7.6", m2, nil, nil, []*KeyRules{Key("F", Skip.When(true), NotNil)}, ""},
		{"t7.7", m2, nil, nil, []*KeyRules{Key("F", Skip.When(false), Required)}, "F: cannot be blank."},
		{"t7.8", m2, nil, nil, []*KeyRules{Key("F", Skip.When(false), NotNil)}, "F: is required."},
		// validatable structs
		{"t8.1", m3, nil, nil, []*KeyRules{Key("M3", Skip)}, ""},
		{"t8.2", m3, nil, nil, []*KeyRules{Key("M3")}, "M3: (A: error abc.)."},
		{"t8.3", m4, nil, nil, []*KeyRules{Key("M3")}, ""},
		// internal error
		{"t9.1", m5, nil, nil, []*KeyRules{Key("A", &validateAbc{}), Key("B", Required), Key("A", &validateInternalError{})}, "error internal"},
		{"t9.2", m5, nil, []Rule{&validateInternalError{}}, nil, "error internal"},
		// shared rules
		{"t10.1", m6, nil, []Rule{&validateXyz{}}, nil, "11: error xyz."},
		{"t10.2", m6, nil, []Rule{&validateXyz{}}, []*KeyRules{Key(11, Length(8, 9)), Key(22, Length(8, 9))}, "11: error xyz; 22: the length must be between 8 and 9."},
		{"t10.3", m7, []Rule{&validateXyz{}}, nil, nil, "abc: error xyz."},
		{"t10.4", m8, []Rule{Length(3, 4)}, nil, []*KeyRules{Key("11", Required), Key("222", &validateXyz{})}, "11: the length must be between 3 and 4."},
	}
	for _, test := range tests {
		err1 := Validate(test.model, Map(test.distinctKeys...).AllowExtraKeys().Keys(test.keyRules...).Values(test.valueRules...))
		err2 := ValidateWithContext(context.Background(), test.model, Map(test.distinctKeys...).AllowExtraKeys().Keys(test.keyRules...).Values(test.valueRules...))
		assertError(t, test.err, err1, test.tag)
		assertError(t, test.err, err2, test.tag)
	}

	a := map[string]interface{}{"Name": "name", "Value": "demo", "Extra": true}
	err := Validate(a, Map(
		Key("Name", Required),
		Key("Value", Required, Length(5, 10)),
	))
	assert.EqualError(t, err, "Extra: key not expected; Value: the length must be between 5 and 10.")
}

func TestMapWithContext(t *testing.T) {
	m1 := map[string]interface{}{"A": "abc", "B": "xyz", "c": "abc", "g": "xyz"}
	m2 := map[string]interface{}{"A": "internal", "B": ""}
	tests := []struct {
		tag   string
		model interface{}
		rules []*KeyRules
		err   string
	}{
		// normal rules
		{"t1.1", m1, []*KeyRules{Key("A", &validateContextAbc{}), Key("B", &validateContextXyz{})}, ""},
		{"t1.2", m1, []*KeyRules{Key("A", &validateContextXyz{}), Key("B", &validateContextAbc{})}, "A: error xyz; B: error abc."},
		{"t1.3", m1, []*KeyRules{Key("A", &validateContextXyz{}), Key("c", &validateContextXyz{})}, "A: error xyz; c: error xyz."},
		{"t1.4", m1, []*KeyRules{Key("g", &validateContextAbc{})}, "g: error abc."},
		// skip rule
		{"t2.1", m1, []*KeyRules{Key("g", Skip, &validateContextAbc{})}, ""},
		{"t2.2", m1, []*KeyRules{Key("g", &validateContextAbc{}, Skip)}, "g: error abc."},
		// internal error
		{"t3.1", m2, []*KeyRules{Key("A", &validateContextAbc{}), Key("B", Required), Key("A", &validateInternalError{})}, "error internal"},
	}
	for _, test := range tests {
		err := ValidateWithContext(context.Background(), test.model, Map(test.rules...).AllowExtraKeys())
		assertError(t, test.err, err, test.tag)
	}

	a := map[string]interface{}{"Name": "name", "Value": "demo", "Extra": true}
	err := ValidateWithContext(context.Background(), a, Map(
		Key("Name", Required),
		Key("Value", Required, Length(5, 10)),
	))
	if assert.NotNil(t, err) {
		assert.Equal(t, "Extra: key not expected; Value: the length must be between 5 and 10.", err.Error())
	}
}
