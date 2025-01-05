package customsource

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strconv"
)

var (
	ErrParse = fmt.Errorf("no match")
	regex    = regexp.MustCompile(`^([0-9]+)_(.*)\.(` + string(Down) + `|` + string(Up) + `)\.(.*)$`)
)

type partialDriver struct {
	migrations *migrations
	fsys       fs.FS
	path       string
}

func (d *partialDriver) init(fsys fs.FS, path string) error {
	entries, err := fs.ReadDir(fsys, path)
	if err != nil {
		return err
	}

	ms := newMigrations()
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		m, err := parse(e.Name())
		if err != nil {
			continue
		}

		file, err := e.Info()
		if err != nil {
			return err
		}

		if !ms.Append(m) {
			return errDuplicateMigration{
				migration: *m,
				FileInfo:  file,
			}
		}
	}

	d.fsys = fsys
	d.path = path
	d.migrations = ms
	return nil
}

func parse(raw string) (*migration, error) {
	m := regex.FindStringSubmatch(raw)
	if len(m) == 5 {
		versionUint64, err := strconv.ParseUint(m[1], 10, 64)
		if err != nil {
			return nil, err
		}

		return &migration{
			version:    uint(versionUint64),
			identifier: m[2],
			direction:  direction(m[3]),
			raw:        raw,
		}, nil
	}

	return nil, ErrParse
}

type errDuplicateMigration struct {
	migration
	os.FileInfo
}

func (e errDuplicateMigration) Error() string {
	return "duplicate migration file: " + e.Name()
}
