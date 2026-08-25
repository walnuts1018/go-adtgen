package multi

import (
	"encoding/json"
	"testing"
)

func TestGeneratedOutputsBuildAndBehave(t *testing.T) {
	alpha := Alpha(&Left{Name: "left"})
	left, ok := AsAlphaLeft(alpha)
	if !ok {
		t.Fatal("AsAlphaLeft(alpha) ok = false, want true")
	}
	if left.Name != "left" {
		t.Fatalf("AsAlphaLeft(alpha).Name = %q, want %q", left.Name, "left")
	}
	if _, ok := AsAlphaRight(alpha); ok {
		t.Fatal("AsAlphaRight(alpha) ok = true, want false")
	}
	got := MatchAlpha(alpha,
		func(v Left) string { return v.Name },
		func(v Right) string { return "unexpected" },
	)
	if got != "left" {
		t.Fatalf("MatchAlpha() = %q, want %q", got, "left")
	}

	beta := NewBeta(Primary{ID: "id-1"}, Secondary{Enabled: true})
	if beta.ID != "id-1" {
		t.Fatalf("NewBeta().ID = %q, want %q", beta.ID, "id-1")
	}
	if !beta.Enabled {
		t.Fatal("NewBeta().Enabled = false, want true")
	}
	if primary := beta.ToPrimary(); primary.ID != "id-1" {
		t.Fatalf("Beta.ToPrimary().ID = %q, want %q", primary.ID, "id-1")
	}
	if secondary := beta.ToSecondary(); !secondary.Enabled {
		t.Fatal("Beta.ToSecondary().Enabled = false, want true")
	}

	// Test Gamma (Left belongs to both Alpha and Gamma, Gamma uses discriminator)
	var gamma Gamma = &Left{Name: "shared-left"}
	leftFromGamma, ok := AsGammaLeft(gamma)
	if !ok || leftFromGamma.Name != "shared-left" {
		t.Fatalf("AsGammaLeft(gamma) = (%+v, %t), want shared-left", leftFromGamma, ok)
	}

	gammaJSON, err := MarshalGamma(gamma)
	if err != nil {
		t.Fatalf("MarshalGamma() error = %v", err)
	}
	var rawMap map[string]any
	if err := json.Unmarshal(gammaJSON, &rawMap); err != nil {
		t.Fatalf("json.Unmarshal(gammaJSON) error = %v", err)
	}
	if rawMap["event_type"] != "Left" || rawMap["Name"] != "shared-left" {
		t.Fatalf("gammaJSON content = %+v, want event_type=Left Name=shared-left", rawMap)
	}

	unmarshaledGamma, err := UnmarshalGamma(gammaJSON)
	if err != nil {
		t.Fatalf("UnmarshalGamma() error = %v", err)
	}
	if leftG, ok := AsGammaLeft(unmarshaledGamma); !ok || leftG.Name != "shared-left" {
		t.Fatalf("AsGammaLeft(unmarshaledGamma) = (%+v, %t), want shared-left", leftG, ok)
	}
}
