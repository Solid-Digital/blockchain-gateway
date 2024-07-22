package ares

import (
	"io"
)

type FileStore interface {
	GetObject(id string) (io.ReadCloser, error)
	PutObject(id string, reader io.Reader) error
}
