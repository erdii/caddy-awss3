/**
 * Attribution:
 * This file contains code taken from coopernurse/caddy-awslambda.
 */

package awss3

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/mholt/caddy/caddyhttp/httpserver"
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

// MustNewRequest creates a new http request with caddy context (unused for now)
func MustNewRequest(method, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		panic(err)
	}
	replacer := httpserver.NewReplacer(req, nil, "")
	newContext := context.WithValue(req.Context(), httpserver.ReplacerCtxKey, replacer)
	req = req.WithContext(newContext)
	return req
}
