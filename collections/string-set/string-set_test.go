package ss

import "testing"

func TestStringSliceAdd(t *testing.T) {
	ss := NewStringSet()
	ss.Add("value")

	if ss.Len() != 1 {
		t.Fatalf("expected length to be 1")
	}

	if ss.Slice()[0] != "value" {
		t.Fatalf("expected value to be 'value'")
	}
}

func TestStringSliceAddMultiple(t *testing.T) {
	ss := NewStringSet()
	ss.Add("v2")
	t.Logf("slice after add: %q", ss.slice)

	ss.Add("v1")
	t.Logf("slice after add: %q", ss.slice)

	if ss.Len() != 2 {
		t.Fatalf("expected length to be 2; actual length is %d", ss.Len())
	}

	if ss.Slice()[0] != "v1" || ss.Slice()[1] != "v2" {
		t.Fatalf("StringSet does not contain expected values; actual = %q", ss.slice)
	}
}

func TestStringSliceAddDuplicate(t *testing.T) {
	ss := NewStringSet()
	ss.Add("v2")
	ss.Add("v1")
	ss.Add("v2")

	if ss.Len() != 2 {
		t.Fatalf("expected length to be 2; actual length is %d", ss.Len())
	}

	if ss.Slice()[0] != "v1" || ss.Slice()[1] != "v2" {
		t.Fatalf("StringSet does not contain expected values; actual = %q", ss.slice)
	}
}

func TestStringSliceRemove(t *testing.T) {
	ss := NewStringSet()
	ss.Add("v3")
	ss.Add("v2")
	ss.Add("v1")

	ss.Remove("v2")
	t.Logf("Slice after remove: %q", ss.slice)

	if ss.Len() != 2 {
		t.Fatalf("expected length to be 2; actual length is %d", ss.Len())
	}

	if ss.Slice()[0] != "v1" || ss.Slice()[1] != "v3" {
		t.Fatalf("StringSet does not contain expected values; actual = %q", ss.slice)
	}
}

func TestStringSliceRemoveNonExistant(t *testing.T) {
	ss := NewStringSet()
	ss.Add("v2")
	ss.Add("v1")

	ss.Remove("v4")

	if ss.Len() != 2 {
		t.Fatalf("expected length to be 2; actual length is %d", ss.Len())
	}

	if ss.Slice()[0] != "v1" || ss.Slice()[1] != "v2" {
		t.Fatalf("StringSet does not contain expected values; actual = %q", ss.slice)
	}
}

func TestStringSliceContains(t *testing.T) {
	ss := NewStringSet()
	ss.Add("v3")
	ss.Add("v2")
	ss.Add("v1")

	if !ss.Contains("v1") || !ss.Contains("v2") || !ss.Contains("v3") || ss.Contains("v4") {
		t.Fatalf("StringSet does not contain expected values; actual = %q", ss.slice)
	}

	if ss.Len() != 3 {
		t.Fatalf("StringSet does not contain expected values; actual = %q", ss.slice)
	}
}
