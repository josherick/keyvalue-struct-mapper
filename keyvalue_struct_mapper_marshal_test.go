package struct_mapper

import "testing"

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
		name     string
		key      string
		value    string
		rawValue interface{}
	}{
		{
			name:     "string marshal",
			key:      "test_string",
			value:    "mystring",
			rawValue: "mystring",
		},
		{
			name:     "bool marshal 'true'",
			key:      "test-bool",
			value:    "true",
			rawValue: true,
		},
		{
			name:     "bool marshal 'false'",
			key:      "test-bool-false",
			value:    "false",
			rawValue: false,
		},
		{
			name:     "bool marshal '1'",
			key:      "test-bool-with-one",
			value:    "true",
			rawValue: true,
		},
		{
			name:     "bool marshal '0'",
			key:      "test-bool-with-zero",
			value:    "false",
			rawValue: false,
		},
		{
			name:     "float marshal",
			key:      "test-float64",
			value:    "3.141592",
			rawValue: 3.141592,
		},
		{
			name:     "int inside struct marshal",
			key:      "test_int",
			value:    "42",
			rawValue: 42,
		},
		{
			name:     "int with substitution marshal",
			key:      "test_int:1234:with_substitution",
			value:    "84",
			rawValue: 84,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &staticStructMapper{
				"1234",
				make(map[string]string),
				make(map[string]interface{}),
			}

			// Sets string version.
			New(nil, s, nil, s).Marshal(&c)
			val, ok := s.setValues[test.key]
			if !ok {
				t.Fatalf("string value was not set for key %+v", test.key)
			}
			if val != test.value {
				t.Fatalf("string value was not set properly, expected: %+v, received %+v", test.value, val)
			}

			// Sets raw version.
			New(nil, nil, s, s).Marshal(&c)
			rawVal, ok := s.setRawValues[test.key]
			if !ok {
				t.Fatalf("raw value was not set for key %+v", test.key)
			}
			if rawVal != test.rawValue {
				t.Fatalf("raw value was not set properly, expected: %+v, received %+v", test.value, rawVal)
			}

			// Sets both.
			s.setValues = make(map[string]string)
			s.setRawValues = make(map[string]interface{})
			New(nil, s, s, s).Marshal(&c)
			strVal, strOk := s.setValues[test.key]
			rawVal, rawOk := s.setRawValues[test.key]
			if !strOk || !rawOk {
				t.Fatalf("both value was not set for key %+v", test.key)
			}
			if strVal != test.value || rawVal != test.rawValue {
				t.Fatalf("both value was not set properly, expected: %+v, received %+v", test.value, strVal)
			}
		})
	}
}
