package adapter

import "testing"

func TestNewTaggedMessage(t *testing.T) {
	if tag := NewTaggedMessage(nil).Tag; tag != 1 {
		t.Fatalf("expected tag %d, got %d", 1, tag)
	}

	if tag := NewTaggedMessage(nil).Tag; tag != 2 {
		t.Fatalf("expected tag %d, got %d", 2, tag)
	}

	if tag := NewTaggedMessage(nil).Tag; tag != 3 {
		t.Fatalf("expected tag %d, got %d", 3, tag)
	}
}
