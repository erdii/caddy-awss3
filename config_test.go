package awss3

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

// Taken from coopernurse/caddy-awslambda
func TestToAwsConfigStaticCreds(t *testing.T) {
	c := &Config{
		AwsAccess: "a-key",
		AwsSecret: "secret",
	}
	expected := credentials.NewStaticCredentials("a-key", "secret", "")
	actual := c.ToAwsConfig()
	if !reflect.DeepEqual(expected, actual.Credentials) {
		t.Errorf("\nExpected: %v\n  Actual: %v", expected, actual.Credentials)
	}
}

// Taken from coopernurse/caddy-awslambda
func TestToAwsConfigStaticRegion(t *testing.T) {
	c := &Config{
		AwsRegion: "us-west-2",
	}
	expected := aws.NewConfig()
	actual := c.ToAwsConfig()
	if c.AwsRegion != *actual.Region {
		t.Errorf("\nExpected: %v\n  Actual: %v", c.AwsRegion, *actual.Region)
	}
	if !reflect.DeepEqual(expected.Credentials, actual.Credentials) {
		t.Errorf("\nExpected: %v\n  Actual: %v", expected.Credentials, actual.Credentials)
	}
}

// Taken from coopernurse/caddy-awslambda
func TestToAwsConfigDefaults(t *testing.T) {
	c := &Config{}
	expected := aws.NewConfig()
	actual := c.ToAwsConfig()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nExpected: %v\n  Actual: %v", expected, actual.Credentials)
	}
}

// Taken from coopernurse/caddy-awslambda
func TestParseConfigs(t *testing.T) {
	for i, test := range []struct {
		input    string
		expected []*Config
	}{
		{"awss3 test-bucket", []*Config{&Config{
			Bucket: "test-bucket",
		}}},
		{`awss3 another-bucket {
    aws_access my-access
    aws_secret my-secret
    aws_region eu-central-1
}`,
			[]*Config{
				&Config{
					Bucket:    "another-bucket",
					AwsAccess: "my-access",
					AwsSecret: "my-secret",
					AwsRegion: "eu-central-1",
				},
			},
		},
		{`awss3 first-bucket {
    aws_region eu-west-1
}
awss3 second-bucket {
    aws_region us-east-1
}`,
			[]*Config{
				&Config{
					Bucket:    "first-bucket",
					AwsRegion: "eu-west-1",
				},
				&Config{
					Bucket:    "second-bucket",
					AwsRegion: "us-east-1",
				},
			},
		},
	} {
		controller := caddy.NewTestController("http", test.input)
		actual, err := ParseConfigs(controller)
		if err != nil {
			t.Errorf("ParseConfigs return err: %v", err)
		}
		// for i := range actual {
		// 	actual[i].invoker = nil
		// }
		EqOrErr(test.expected, actual, i, t)
	}
}

// Taken from coopernurse/caddy-awslambda
func mustNewRequest(method, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		panic(err)
	}
	replacer := httpserver.NewReplacer(req, nil, "")
	newContext := context.WithValue(req.Context(), httpserver.ReplacerCtxKey, replacer)
	req = req.WithContext(newContext)
	return req
}
