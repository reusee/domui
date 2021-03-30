package domui

import "testing"

func TestSortedMap(t *testing.T) {
	m := new(SortedMap)
	m.Set("1", true)
	m.Set("3", true)
	m.Set("2", true)

	v, ok := m.Get("3")
	if !ok {
		t.Fatal()
	}
	if !v.(bool) {
		t.Fatal()
	}

	_, ok = m.Get("4")
	if ok {
		t.Fatal()
	}

	m.Set("1", false)
	v, ok = m.Get("1")
	if !ok {
		t.Fatal()
	}
	if v.(bool) {
		t.Fatal()
	}

}
