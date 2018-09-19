package awss3

import "net/http"

// WriteAtWrapper wraps an http.ResponseWriter instace so one can (unsafely) writeAt into it
type WriteAtWrapper struct {
	w *http.ResponseWriter
}

// NewWriteAtWrapper creates a new WriteAtWrapper instance
func NewWriteAtWrapper(w *http.ResponseWriter) WriteAtWrapper {
	return WriteAtWrapper{
		w: w,
	}
}

// WriteAt exposes an interface to write into the "contained" http.ResponseWriter
// NEVER EVER use WriteAt with non consecutive chunks
func (wrp *WriteAtWrapper) WriteAt(p []byte, _pos int64) (int, error) {
	return (*(*wrp).w).Write(p)
}
