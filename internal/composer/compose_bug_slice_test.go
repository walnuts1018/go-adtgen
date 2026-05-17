package composer

import (
	"testing"
)

func TestBaseInputNameSlice(t *testing.T) {
	name := baseInputName("[]foo.Bar")
	if name == "" {
		t.Fatal("empty name")
	}
}
