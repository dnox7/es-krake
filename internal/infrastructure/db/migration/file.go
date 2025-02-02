package migration

import (
	"net/url"
	"os"
	"path/filepath"
)

type file struct {
	partialDriver
	url  string
	path string
}

func (f *file) Open(urlStr string) (*file, error) {
	p, err := parseUrl(urlStr)
	if err != nil {
		return nil, err
	}

	nf := &file{
		url:  urlStr,
		path: p,
	}

	if err := nf.init(os.DirFS(p), "."); err != nil {
		return nil, err
	}
	return nf, nil
}

func (f *file) GetLastestIndex() uint {
	return f.migrations.index[len(f.migrations.index)-1]
}

func parseUrl(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	p := u.Opaque
	if len(p) == 0 {
		p = u.Host + u.Path
	}

	if len(p) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		p = wd
	} else if p[0:1] == "." || p[0:1] != "/" {
		abs, err := filepath.Abs(p)
		if err != nil {
			return "", err
		}
		p = abs
	}
	return p, nil
}
