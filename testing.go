package awss3

import (
	"encoding/json"
	"reflect"
	"testing"
)

// EqOrErr checks for deep equality - taken from coopernurse/caddy-awslambda
func EqOrErr(expected, actual interface{}, num int, t *testing.T) bool {
	if !reflect.DeepEqual(expected, actual) {
		ex, err := json.Marshal(expected)
		ac, err2 := json.Marshal(actual)
		if err != nil || err2 != nil {
			t.Errorf("\nTest %d\nExpected: %+v\n  Actual: %+v", num, expected, actual)
			return false
		}
		t.Errorf("\nTest %d\nExpected: %s\n  Actual: %s", num, ex, ac)
		return false
	}
	return true
}
