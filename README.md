# Key/Value Struct Mapper

Based on [Kelsey Hightower's
envconfig](https://github.com/kelseyhightower/envconfig), though doesn't
support all the features of that package.

Allows easy mapping of key/value data (e.g. from a database) into a typed golang
struct. When data comes back from the store untyped (i.e. as a string), saves a
lot of boilerplate.

Data can also be "repacked" from the struct into your database.

Especially useful with Redis, but there's a bunch of different uses.

Not entirely sure if "marshal"/"unmarshal" are the technically correct terms
here, but seemed close enough.

Supports the types in the example below (and tests). Other types have not been
tested.

Example Usage:
```
import (
	"fmt"
	"strings"
	struct_mapper "github.com/josherick/keyvalue-struct-mapper"
)

// Type to "unpack" the data into.
type config struct {
	TestString      string  `keyname:"test_string"`
	TestBool        bool    `keyname:"test-bool"`
	TestBoolWithOne bool    `keyname:"test-bool-with-one"`
	TestFloat64     float64 `keyname:"test-float64"`
	TestStruct      struct {
		TestInt                 int `keyname:"test_int"`
		TestIntWithSubstitution int `keyname:"test_int:%s:with_substitution"`
	} `keyname:"should-not-be-sent"`
}

// Data provider and consumer for mapper.
type staticStructMapper string

func (s staticStructMapper) ProcessKey(key string) string {
	return strings.Replace(key, "%s", string(s), -1)
}

// Called for each key in the struct after calling Unmarshal.
// Returns a value from data store (in this case, stati
func (staticStructMapper) Get(key string) (string, bool) {
	store := map[string]string{
		"test_string":                     "mystring",
		"test-bool":                       "false",
		"test-bool-with-one":              "1",
		"test-float64":                    "3.141592",
		"test_int":                        "42",
		"test_int:1234:with_substitution": "84",
	}
	val, ok := store[key]
	return val, ok
}

// Called for each key in the struct after calling Marshal.
// Can set data store here with serialized values from the struct.
func (s staticStructMapper) Set(key string, value string) {
	fmt.Printf("Setting %v to key %s\n", key, value)
}

func main() {
	s := staticStructMapper("1234")
	var c config

	// Populate c with values from our data store (a map, `store`)
	struct_mapper.New(s, s, s).Unmarshal(&c)
	fmt.Printf("Unpacked into c: %+v\n", c)

	// Will print each time we set a value. We could instead put this back into
	// the store.
	struct_mapper.New(s, s, s).Marshal(&c)
}
```

Output:
```
Unpacked into c: {TestString:mystring TestBool:false TestBoolWithOne:true TestFloat64:3.141592 TestStruct:{TestInt:42 TestIntWithSubstitution:84}}
Setting test_string to key mystring
Setting test-bool to key false
Setting test-bool-with-one to key true
Setting test-float64 to key 3.141592
Setting test_int to key 42
Setting test_int:1234:with_substitution to key 84
```
