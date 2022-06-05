package struct_mapper

import "testing"

func TestUnmarshal(t *testing.T) {
	s := &staticStructMapper{
		"1234",
		make(map[string]string),
		make(map[string]interface{}),
	}

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
			New(s, nil, nil, s).Unmarshal(&c)
			if test.getValue(c) != test.correctValue {
				t.Fatalf("value was not set properly, config: %+v", c)
			}
		})
	}
}
