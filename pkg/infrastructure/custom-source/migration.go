package customsource

import "sort"

type (
	direction string
	uintSlice []uint
)

const (
	Down direction = "down"
	Up   direction = "up"
)

type migration struct {
	version    uint
	identifier string
	direction  direction
	raw        string
}

type migrations struct {
	index      uintSlice
	migrations map[uint]map[direction]*migration
}

func newMigrations() *migrations {
	return &migrations{
		index:      make(uintSlice, 0),
		migrations: make(map[uint]map[direction]*migration),
	}
}

func (m *migrations) Append(ele *migration) (ok bool) {
	if ele == nil {
		return false
	}

	if m.migrations[ele.version] == nil {
		m.migrations[ele.version] = make(map[direction]*migration)
	}

	if _, dup := m.migrations[ele.version][ele.direction]; dup {
		return false
	}

	m.migrations[ele.version][ele.direction] = ele
	m.buildIndex()

	return true
}

func (m *migrations) buildIndex() {
	m.index = make(uintSlice, 0, len(m.migrations))
	for ver := range m.migrations {
		m.index = append(m.index, ver)
	}
	sort.Slice(m.index, func(i, j int) bool {
		return m.index[i] < m.index[j]
	})
}
