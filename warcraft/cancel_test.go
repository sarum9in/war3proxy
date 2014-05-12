package warcraft

import "testing"

func TestCancel(t *testing.T) {
    Cancel := NewCancel(11)
    a2, err := ParseCancel(Cancel.Bytes())
    if err != nil {
        t.Fatalf("Error: %v", err)
    }
    if a2 != Cancel {
        t.Fatalf("a2 != Cancel")
    }
}
