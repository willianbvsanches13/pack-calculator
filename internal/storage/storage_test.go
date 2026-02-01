package storage

import (
	"reflect"
	"sort"
	"testing"
)

func TestNewMemoryStorage(t *testing.T) {
	s := NewMemoryStorage()

	sizes := s.GetPackSizes()
	expected := DefaultPackSizes()

	sort.Ints(sizes)
	sort.Ints(expected)

	if !reflect.DeepEqual(sizes, expected) {
		t.Errorf("expected %v, got %v", expected, sizes)
	}
}

func TestNewMemoryStorageWithSizes(t *testing.T) {
	custom := []int{100, 200, 300}
	s := NewMemoryStorageWithSizes(custom)

	sizes := s.GetPackSizes()
	if !reflect.DeepEqual(sizes, custom) {
		t.Errorf("expected %v, got %v", custom, sizes)
	}
}

func TestGetPackSizesReturnsCopy(t *testing.T) {
	s := NewMemoryStorage()

	sizes := s.GetPackSizes()
	original := sizes[0]
	sizes[0] = 99999

	sizes2 := s.GetPackSizes()
	if sizes2[0] != original {
		t.Error("GetPackSizes should return a copy, not the original slice")
	}
}

func TestSetPackSizes(t *testing.T) {
	s := NewMemoryStorage()

	newSizes := []int{10, 20, 30}
	s.SetPackSizes(newSizes)

	sizes := s.GetPackSizes()
	if !reflect.DeepEqual(sizes, newSizes) {
		t.Errorf("expected %v, got %v", newSizes, sizes)
	}
}

func TestSetPackSizesMakesCopy(t *testing.T) {
	s := NewMemoryStorage()

	newSizes := []int{10, 20, 30}
	s.SetPackSizes(newSizes)

	newSizes[0] = 99999

	sizes := s.GetPackSizes()
	if sizes[0] == 99999 {
		t.Error("SetPackSizes should make a copy of the input slice")
	}
}

func TestAddPackSize(t *testing.T) {
	s := NewMemoryStorageWithSizes([]int{100, 200})

	if !s.AddPackSize(300) {
		t.Error("AddPackSize should return true for new size")
	}

	sizes := s.GetPackSizes()
	found := false
	for _, size := range sizes {
		if size == 300 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 300 to be added")
	}
}

func TestAddPackSizeDuplicate(t *testing.T) {
	s := NewMemoryStorageWithSizes([]int{100, 200})

	if s.AddPackSize(100) {
		t.Error("AddPackSize should return false for duplicate")
	}

	sizes := s.GetPackSizes()
	if len(sizes) != 2 {
		t.Errorf("expected 2 sizes, got %d", len(sizes))
	}
}

func TestRemovePackSize(t *testing.T) {
	s := NewMemoryStorageWithSizes([]int{100, 200, 300})

	if !s.RemovePackSize(200) {
		t.Error("RemovePackSize should return true for existing size")
	}

	sizes := s.GetPackSizes()
	for _, size := range sizes {
		if size == 200 {
			t.Error("expected 200 to be removed")
		}
	}

	if len(sizes) != 2 {
		t.Errorf("expected 2 sizes, got %d", len(sizes))
	}
}

func TestRemovePackSizeNotFound(t *testing.T) {
	s := NewMemoryStorageWithSizes([]int{100, 200})

	if s.RemovePackSize(999) {
		t.Error("RemovePackSize should return false for non-existing size")
	}

	sizes := s.GetPackSizes()
	if len(sizes) != 2 {
		t.Errorf("expected 2 sizes, got %d", len(sizes))
	}
}

func TestDefaultPackSizes(t *testing.T) {
	defaults := DefaultPackSizes()

	if len(defaults) != 5 {
		t.Errorf("expected 5 default pack sizes, got %d", len(defaults))
	}

	expected := []int{250, 500, 1000, 2000, 5000}
	sort.Ints(defaults)
	sort.Ints(expected)

	if !reflect.DeepEqual(defaults, expected) {
		t.Errorf("expected %v, got %v", expected, defaults)
	}
}
