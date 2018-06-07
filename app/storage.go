package app

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"net/url"
)

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

// NopCloser returns a ReadCloser with a no-op Close method wrapping
// the provided Reader r.
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w}
}

func DummyStorage() Storage {
	return &dummyStorage{}
}

type dummyStorage struct {
}

type Storage interface {
	Create(u *url.URL) (io.WriteCloser, error)
}

func (s *dummyStorage) Create(u *url.URL) (io.WriteCloser, error) {
	return NopWriteCloser(ioutil.Discard), nil
}

func FileStorage(dst string) Storage {
	return &fileStorage{
		dst: dst,
	}
}

type fileStorage struct {
	dst string
}

func (s *fileStorage) Create(u *url.URL) (io.WriteCloser, error) {
	parts := strings.Split(u.Path, "/")

	dst := []string{
		s.dst,
	}

	dst = append(dst, u.Host)
	dst = append(dst, parts...)

	fileName := "index.html"

	if path.Ext(u.Path) != "" {
		// dir
		fileName = dst[len(dst)-1]
		dst = dst[0 : len(dst)-1]
	}

	if err := os.MkdirAll(path.Join(dst...), 0700); err != nil {
		return nil, err
	}

	f, err := os.Create(path.Join(path.Join(dst...), fileName))
	if err != nil {
		return nil, err
	}

	return f, nil
}
