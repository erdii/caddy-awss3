package awss3

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mholt/caddy"
)

// Config specifies configuration for a single awslambda block
type Config struct {
	// Request base path
	Path string
	// AWS S3 Bucket this config block maps to
	Bucket string
	// AWS Access Key. If omitted, AWS_ACCESS_KEY_ID env var is used.
	AwsAccess string
	// AWS Secret Key. If omitted, AWS_SECRET_ACCESS_KEY env var is used.
	AwsSecret string
	// AWS Region. If omitted, AWS_REGION env var is used.
	AwsRegion string
	// AWS S3 client instance
	S3Client *s3.S3
}

// ToAwsConfig returns a new *aws.Config instance using the AWS related values on Config.
// If AwsRegion is empty, the AWS_REGION env var is used.
// If AwsAccess is empty, the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY env vars are used.
// Taken from coopernurse/caddy-awslambda
func (c *Config) ToAwsConfig() *aws.Config {
	awsConf := aws.NewConfig()
	if c.AwsRegion != "" {
		awsConf.WithRegion(c.AwsRegion)
	}
	if c.AwsAccess != "" {
		awsConf.WithCredentials(credentials.NewStaticCredentials(
			c.AwsAccess, c.AwsSecret, "",
		))
	}
	return awsConf
}

// InitS3Client creates a new s3 client instance for a specific bucket
func (c *Config) initS3Client() error {
	sess, err := session.NewSession(c.ToAwsConfig())
	if err != nil {
		return err
	}
	c.S3Client = s3.New(sess)
	return nil
}

// ParseConfigs parses a Caddy awslambda config block into a Config struct.
// Inspired by coopernurse/caddy-awslambda
func ParseConfigs(c *caddy.Controller) ([]*Config, error) {
	var configs []*Config
	var conf *Config
	last := ""

	for c.Next() {
		val := c.Val()
		lastTmp := last
		last = ""
		switch lastTmp {
		case "awss3":
			s := []string{val}
			s = append(s, c.RemainingArgs()...)
			switch len(s) {
			case 1:
				conf = &Config{
					Path: s[0],
				}
			case 2:
				conf = &Config{
					Path:   s[0],
					Bucket: s[1],
				}
			default:
				return nil, errors.New("bare awss3 directive usage: awss3 path [bucket]")
			}
			configs = append(configs, conf)
		case "bucket":
			conf.Bucket = val
		case "aws_access":
			conf.AwsAccess = val
		case "aws_secret":
			conf.AwsSecret = val
		case "aws_region":
			conf.AwsRegion = val
		default:
			last = val
		}
	}

	for _, conf := range configs {
		err := conf.initS3Client()
		if err != nil {
			return nil, err
		}
	}

	return configs, nil
}

// StripPathPrefix strips the basepath from a request path
// Taken from coopernurse/caddy-awslambda
func (c *Config) StripPathPrefix(reqPath string) string {
	prefix := c.Path

	if strings.HasPrefix(reqPath, prefix) {
		reqPath = reqPath[len(prefix):]
		if !strings.HasPrefix(reqPath, "/") {
			reqPath = "/" + reqPath
		}
	}
	return reqPath
}
