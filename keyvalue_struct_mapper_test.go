package struct_mapper

import (
	"fmt"
	"strings"
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
	setRawValues           map[string]interface{}
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

func (s staticStructMapper) Set(key string, value string) error {
	if val, ok := s.setValues[key]; ok {
		panic(fmt.Sprintf("key %s was already set to %s. Tried to set to %s.", key, val, value))
	}
	s.setValues[key] = value
	return nil
}

func (s staticStructMapper) SetRaw(key string, value interface{}) error {
	if val, ok := s.setRawValues[key]; ok {
		panic(fmt.Sprintf("key %s was already set to %s. Tried to set to %s.", key, val, value))
	}
	s.setRawValues[key] = value
	return nil
}
