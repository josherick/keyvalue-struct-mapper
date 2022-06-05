package struct_mapper

import (
	"fmt"
	"strings"
	"testing"
)

type config struct {
	TestString       string  `keyname:"test_string"`
	TestBool         bool    `keyname:"test-bool"`
	TestBoolFalse    bool    `keyname:"test-bool-false"`
	TestBoolWithOne  bool    `keyname:"test-bool-with-one"`
	TestBoolWithZero bool    `keyname:"test-bool-with-zero"`
	TestFloat64      float64 `keyname:"test-float64"`
	TestStruct       struct {
		TestInt                 int `keyname:"test_int"`
		TestIntWithSubstitution int `keyname:"test_int:%s:with_substitution"`
	} `keyname:"should-not-be-sent"`
}

// Implements KeyValueGetter, KeyProcessor, and KeyValueSetter.
// - Returns static values.
// - Maps "%s" in keys to the value of replacementForPercentS.
// - When a value is set, puts it into setValues. If the value was already set,
// panics.
type staticStructMapper struct {
	replacementForPercentS string
	setValues              map[string]string
}

func (s staticStructMapper) ProcessKey(key string) string {
	if key == "should-not-be-sent" {
		panic(fmt.Sprintf("Received key from a struct in ProcessKey, this should no be sent"))
	}
	return strings.Replace(key, "%s", s.replacementForPercentS, -1)
}

func (staticStructMapper) Get(key string) (string, bool) {
	if key == "should-not-be-sent" {
		panic(fmt.Sprintf("Received key from a struct in Get, this should no be sent"))
	}
	switch key {
	case "test_string":
		return "mystring", true
	case "test-bool":
		return "true", true
	case "test-bool-false":
		return "false", true
	case "test-bool-with-one":
		return "1", true
	case "test-bool-with-zero":
		return "0", true
	case "test-float64":
		return "3.141592", true
	case "test_int":
		return "42", true
	case "test_int:1234:with_substitution":
		return "84", true
	}
	fmt.Printf("received unknown key: %s\n", key)
	return "", false
}

func (s staticStructMapper) Set(key string, value string) {
	if val, ok := s.setValues[key]; ok {
		panic(fmt.Sprintf("key %s was already set to %s. Tried to set to %s.", key, val, value))
	}
	s.setValues[key] = value
}

func TestUnmarshal(t *testing.T) {
	s := &staticStructMapper{"1234", make(map[string]string)}

	tests := []struct {
		name         string
		checkBefore  bool
		getValue     func(config) interface{}
		correctValue interface{}
	}{
		{
			name:         "string unmarshal",
			checkBefore:  true,
			getValue:     func(c config) interface{} { return c.TestString },
			correctValue: "mystring",
		},
		{
			name:         "bool unmarshal 'true'",
			checkBefore:  true,
			getValue:     func(c config) interface{} { return c.TestBool },
			correctValue: true,
		},
		{
			name:         "bool unmarshal 'false'",
			checkBefore:  false,
			getValue:     func(c config) interface{} { return c.TestBoolFalse },
			correctValue: false,
		},
		{
			name:         "bool unmarshal '1'",
			checkBefore:  true,
			getValue:     func(c config) interface{} { return c.TestBoolWithOne },
			correctValue: true,
		},
		{
			name:         "bool unmarshal '0'",
			checkBefore:  false,
			getValue:     func(c config) interface{} { return c.TestBoolWithZero },
			correctValue: false,
		},
		{
			name:         "float unmarshal",
			checkBefore:  true,
			getValue:     func(c config) interface{} { return c.TestFloat64 },
			correctValue: 3.141592,
		},
		{
			name:         "int inside struct unmarshal",
			checkBefore:  true,
			getValue:     func(c config) interface{} { return c.TestStruct.TestInt },
			correctValue: 42,
		},
		{
			name:         "int with substitution unmarshal",
			checkBefore:  true,
			getValue:     func(c config) interface{} { return c.TestStruct.TestIntWithSubstitution },
			correctValue: 84,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var c config
			if test.checkBefore && test.getValue(c) == test.correctValue {
				t.Fatalf("value was already set before unmarshaling")
			}
			New(s, s, s).Unmarshal(&c)
			if test.getValue(c) != test.correctValue {
				t.Fatalf("value was not set properly, config: %+v", c)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	c := config{
		TestString:       "mystring",
		TestBool:         true,
		TestBoolFalse:    false,
		TestBoolWithOne:  true,
		TestBoolWithZero: false,
		TestFloat64:      3.141592,
	}
	c.TestStruct.TestInt = 42
	c.TestStruct.TestIntWithSubstitution = 84

	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "string marshal",
			key:   "test_string",
			value: "mystring",
		},
		{
			name:  "bool marshal 'true'",
			key:   "test-bool",
			value: "true",
		},
		{
			name:  "bool marshal 'false'",
			key:   "test-bool-false",
			value: "false",
		},
		{
			name:  "bool marshal '1'",
			key:   "test-bool-with-one",
			value: "true",
		},
		{
			name:  "bool marshal '0'",
			key:   "test-bool-with-zero",
			value: "false",
		},
		{
			name:  "float marshal",
			key:   "test-float64",
			value: "3.141592",
		},
		{
			name:  "int inside struct marshal",
			key:   "test_int",
			value: "42",
		},
		{
			name:  "int with substitution marshal",
			key:   "test_int:1234:with_substitution",
			value: "84",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &staticStructMapper{"1234", make(map[string]string)}
			New(s, s, s).Marshal(&c)
			val, ok := s.setValues[test.key]
			if !ok {
				t.Fatalf("value was not set for key %+v", test.key)
			}
			if val != test.value {
				t.Fatalf("value was not set properly, expected: %+v, received %+v", test.value, val)
			}
		})
	}
}
