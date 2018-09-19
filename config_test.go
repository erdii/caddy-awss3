/**
 * Attribution:
 * This file contains code taken from coopernurse/caddy-awslambda.
 */

package awss3

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/mholt/caddy"
)

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

func TestToAwsConfigDefaults(t *testing.T) {
	c := &Config{}
	expected := aws.NewConfig()
	actual := c.ToAwsConfig()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nExpected: %v\n  Actual: %v", expected, actual.Credentials)
	}
}

func TestParseConfigs(t *testing.T) {
	for i, test := range []struct {
		input    string
		expected []*Config
	}{
		{"awss3 / test-bucket", []*Config{&Config{
			Path:   "/",
			Bucket: "test-bucket",
		}}},
		{`awss3 /path another-bucket {
    aws_access my-access
    aws_secret my-secret
    aws_region eu-central-1
}`,
			[]*Config{
				&Config{
					Path:      "/path",
					Bucket:    "another-bucket",
					AwsAccess: "my-access",
					AwsSecret: "my-secret",
					AwsRegion: "eu-central-1",
				},
			},
		},
		{`awss3 /foo first-bucket {
    aws_region eu-west-1
}
awss3 /bar {
    aws_region us-east-1
    bucket second-bucket
}`,
			[]*Config{
				&Config{
					Path:      "/foo",
					Bucket:    "first-bucket",
					AwsRegion: "eu-west-1",
				},
				&Config{
					Path:      "/bar",
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
		for i := range actual {
			actual[i].S3Client = nil
		}
		EqOrErr(test.expected, actual, i, t)
	}
}

func TestStripPathPrefix(t *testing.T) {
	c := Config{
		Path:   "/files/",
		Bucket: "test-bucket",
	}

	for i, test := range []struct {
		reqPath  string
		expected string
	}{
		{"/files/foo", "/foo"},
		{"/files/blahstuff/other/things", "/blahstuff/other/things"},
		{"/files/foo", "/foo"},
		{"/other/foo", "/other/foo"},
		{"/other/bar", "/other/bar"},
	} {
		actual := c.StripPathPrefix(test.reqPath)
		if actual != test.expected {
			t.Errorf("Test %d failed:\nExpected: %s\n  Actual: %s", i, test.expected, actual)
		}
	}
}
