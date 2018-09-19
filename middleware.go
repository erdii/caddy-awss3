/**
 * Attribution:
 * This file contains code taken from coopernurse/caddy-awslambda.
 */

package awss3

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

	matched, c, path, err := h.match(r)

	if !matched {
		return 404, nil
	}

	if err != nil {
		return 0, err
	}

	if path == "/" {
		// bucket listing is not implemented for now
		return 501, nil
	}

	// make an upfront request to s3, so we can find out header data about the requested object
	result, err := c.S3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: &c.Bucket,
		Key:    &path,
	})
	if err != nil {
		// send status 404 if aws has no such object
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "NotFound" {
				return 404, nil
			}
		}
		return 0, err
	}

	// send 404 for now if the requested object is a folder
	if *result.ContentType == "application/x-directory" {
		return 404, nil
	}

	// set response headers
	w.Header().Set("content-type", *result.ContentType)
	w.Header().Set("ETag", *result.ETag)
	w.Header().Set("Last-Modified", (*result.LastModified).UTC().Format(http.TimeFormat))

	// wrap our http response in a WriteAtWrapper, so we can expose a WriteAt(b []byte, pos int64) function to s3manager
	wrappedWriter := NewWriteAtWrapper(&w)

	downloader := s3manager.NewDownloaderWithClient(c.S3Client, func(d *s3manager.Downloader) {
		// set downloader concurrency to 1, so the file gets downloaded sequentially
		d.Concurrency = 1
	})

	// start download and write the response into wrappedWriter
	_, err = downloader.Download(&wrappedWriter, &s3.GetObjectInput{
		Bucket: &c.Bucket,
		Key:    &path,
	})

	if err != nil {
		return 0, err
	}

	return 200, nil
}

// match finds the best match for a proxy config based on r.
func (h Handler) match(r *http.Request) (bool, *Config, string, error) {
	for _, conf := range h.Configs {
		basePath := conf.Path
		if httpserver.Path(r.URL.Path).Matches(basePath) {
			return true, conf, conf.StripPathPrefix(r.URL.Path), nil
		}
	}

	return false, nil, "", nil
}
