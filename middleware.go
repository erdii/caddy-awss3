package awss3

import (
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

// Handler represents a middleware instance that can gateway requests to AWS S3
type Handler struct {
	Next    httpserver.Handler
	Configs []*Config
}

// ServeHTTP satisfies the httpserver.Handler interface by proxying
// the request to AWS S3
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if len(h.Configs) == 0 {
		return h.Next.ServeHTTP(w, r)
	}

	c, path, err := h.match(r)
	if err != nil {
		return 0, err
	}

	if path == "/" {
		// bucket listing is not implemented for now
		return 501, nil
	}

	result, err := c.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		// send status 404 if aws has no such object
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchKey {
				return 404, nil
			}
		}
		return 0, err
	}

	fmt.Printf("result:\n%+v\n", result)

	contentType := *result.ContentType

	// send 404 for now if the requested object is a folder
	if contentType == "application/x-directory" {
		return 404, nil
	}

	// set received content type
	w.Header().Set("content-type", contentType)

	// send s3 response
	_, err = io.Copy(w, result.Body)
	if err != nil {
		return 0, err
	}

	return 200, nil

	// status := 201

	// w.WriteHeader(status)
	// w.Header().Set("content-type", "application/json")

	// bodyBytes, err := json.Marshal(Response{
	// 	This:   "is a test",
	// 	Bucket: h.Configs[0].Bucket,
	// })

	// if err != nil {
	// 	return 0, err
	// }

	// _, err = w.Write(bodyBytes)

	// if err != nil || status >= 400 {
	// 	return 0, err
	// }

	// return status, nil
}

// match finds the best match for a proxy config based on r.
func (h Handler) match(r *http.Request) (*Config, string, error) {
	var c *Config
	var path string

	for _, conf := range h.Configs {
		basePath := conf.Path
		if httpserver.Path(r.URL.Path).Matches(basePath) {
			c = conf
			path = c.StripPathPrefix(r.URL.Path)
		}
	}

	return c, path, nil
}
