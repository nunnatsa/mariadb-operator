package datastructures

import (
	"reflect"
	"sort"
	"testing"
)

func TestDataStructures(t *testing.T) {
	idx := newIndex("a", "b", "b", "c", "d", "e")
	if !reflect.DeepEqual(idx, newIndex("a", "b", "c", "d", "e")) {
		t.Error("expeting index to remove duplicates")
	}

	item, err := Get(idx, "b")
	expectedItem := "b"
	if err != nil {
		t.Errorf("expecting error to not have occurred: %v", err)
	}
	if item != expectedItem {
		t.Errorf("expeting item to be %s, got %s", "a", item)
	}

	keys := Keys(idx)
	expectedKeys := []string{"a", "b", "c", "d", "e"}
	sort.Strings(keys)
	sort.Strings(expectedKeys)
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("expecting keys to be %v, got: %v", expectedKeys, keys)
	}

	exists := AllExists(idx, "a", "b", "c")
	expectedExists := true
	if exists != expectedExists {
		t.Errorf("expecting exists to be %v, got: %v", expectedExists, exists)
	}

	exists = AllExists(idx, "a", "b", "c", "z")
	expectedExists = false
	if exists != expectedExists {
		t.Errorf("expecting exists to be %v, got: %v", expectedExists, exists)
	}

	filteredIdx := Filter(idx, "a", "b", "c")
	expectedFilteredIdx := newIndex("a", "b", "c")
	if !reflect.DeepEqual(filteredIdx, expectedFilteredIdx) {
		t.Errorf("expecting filtered index to be %v, got: %v", expectedFilteredIdx, filteredIdx)
	}
}

func TestDataStructuresDiff(t *testing.T) {
	tests := []struct {
		name     string
		current  Index[string]
		previous Index[string]
		wantDiff DiffResult
	}{
		{
			name:     "no diff",
			current:  newIndex("a", "b"),
			previous: newIndex("a", "b"),
			wantDiff: DiffResult{
				Rest: []string{"a", "b"},
			},
		},
		{
			name:     "added",
			current:  newIndex("a", "b", "c", "d"),
			previous: newIndex("a", "b"),
			wantDiff: DiffResult{
				Added: []string{"c", "d"},
				Rest:  []string{"a", "b"},
			},
		},
		{
			name:     "deleted",
			current:  newIndex(),
			previous: newIndex("a", "b"),
			wantDiff: DiffResult{
				Deleted: []string{"a", "b"},
			},
		},
		{
			name:     "added and deleted",
			current:  newIndex("b", "d", "e", "f"),
			previous: newIndex("a", "b", "c"),
			wantDiff: DiffResult{
				Added:   []string{"d", "e", "f"},
				Deleted: []string{"a", "c"},
				Rest:    []string{"b"},
			},
		},
		{
			name:     "no intersection",
			current:  newIndex("d", "e", "f"),
			previous: newIndex("a", "b", "c"),
			wantDiff: DiffResult{
				Added:   []string{"d", "e", "f"},
				Deleted: []string{"a", "b", "c"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := Diff(tt.current, tt.previous)
			if !reflect.DeepEqual(tt.wantDiff, diff) {
				t.Errorf("expecting config to be:\n%v\ngot:\n%v\n", tt.wantDiff, diff)
			}
		})
	}
}

func newIndex(items ...string) Index[string] {
	return NewIndex[string](items, func(s string) string {
		return s
	})
}
