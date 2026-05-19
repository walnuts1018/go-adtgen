package multi

import "testing"

func TestGeneratedOutputsBuildAndBehave(t *testing.T) {
	alpha := Alpha(&Left{Name: "left"})
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
}
