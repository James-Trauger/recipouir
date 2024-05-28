package utils

import (
	"io"
)

func drainClose(r io.ReadCloser) error {
	io.Copy(io.Discard, r)
	return r.Close()
}
