package sum

import (
	"encoding/json"
	"testing"
)

func TestGeneratedSumHelpersAndJSON(t *testing.T) {
	var value HogeOrFuga = &Hoge{
		Common: Common{ID: "h-1"},
		Name:   "walnut",
	}

	if got := value.GetID(); got != "h-1" {
		t.Fatalf("value.GetID() = %q, want %q", got, "h-1")
	}
	if got := value.String(); got != "walnut" {
		t.Fatalf("value.String() = %q, want %q", got, "walnut")
	}

	if _, ok := value.(interface{ SetID(string) }); ok {
		t.Fatal("value unexpectedly implements SetID(string)")
	}

	if got, ok := value.AsHoge(); !ok || got.ID != "h-1" || got.Name != "walnut" {
		t.Fatalf("value.AsHoge() = (%+v, %t), want ID=h-1 Name=walnut", got, ok)
	}

	matched := MatchHogeOrFuga(value,
		func(h Hoge) string { return h.Name },
		func(f Fuga) string { return f.ID },
	)
	if matched != "walnut" {
		t.Fatalf("MatchHogeOrFuga(...) = %q, want %q", matched, "walnut")
	}

	label, number := MatchHogeOrFuga2(value,
		func(h Hoge) (string, int) { return h.ID, len(h.Name) },
		func(f Fuga) (string, int) { return f.ID, f.Age },
	)
	if label != "h-1" || number != len("walnut") {
		t.Fatalf("MatchHogeOrFuga2(...) = (%q, %d), want (%q, %d)", label, number, "h-1", len("walnut"))
	}

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if got := string(data); got != `{"id":"h-1","name":"walnut"}` {
		t.Fatalf("json.Marshal() = %s, want %s", got, `{"id":"h-1","name":"walnut"}`)
	}

	decoded, err := UnmarshalHogeOrFuga(data)
	if err != nil {
		t.Fatalf("UnmarshalHogeOrFuga() error = %v", err)
	}
	if got, ok := decoded.AsHoge(); !ok || got.ID != "h-1" || got.Name != "walnut" {
		t.Fatalf("decoded.AsHoge() = (%+v, %t), want ID=h-1 Name=walnut", got, ok)
	}
	if got := decoded.String(); got != "walnut" {
		t.Fatalf("decoded.String() = %q, want %q", got, "walnut")
	}
}

func TestGeneratedSumRejectsAmbiguousJSON(t *testing.T) {
	_, err := UnmarshalHogeOrFuga([]byte(`{"id":"shared"}`))
	if err == nil {
		t.Fatal("expected error")
	}
}
