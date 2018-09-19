package awss3

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewWriteAtWrapper(t *testing.T) {
	// c := &Config{
	// 	AwsAccess: "a-key",
	// 	AwsSecret: "secret",
	// }
	// expected := credentials.NewStaticCredentials("a-key", "secret", "")
	// actual := c.ToAwsConfig()
	// if !reflect.DeepEqual(expected, actual.Credentials) {
	// 	t.Errorf("\nExpected: %v\n  Actual: %v", expected, actual.Credentials)
	// }
	for _, test := range []struct {
		chunks   [][]byte
		expected []byte
	}{
		{[][]byte{[]byte{255, 255, 255, 0, 0, 127, 0, 192}}, []byte{255, 255, 255, 0, 0, 127, 0, 192}},
	} {
		rr := httptest.NewRecorder()
		casted := http.ResponseWriter(rr)
		wrapper := NewWriteAtWrapper(&casted)

		for i, chunk := range test.chunks {
			wrapper.WriteAt(chunk, int64(i*8))
		}

		result := rr.Result()

		actual, err := ioutil.ReadAll(result.Body)

		if err != nil {
			t.Error(err)
		}

		EqOrErr(test.expected, actual, 1, t)
	}
}
