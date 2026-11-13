package slug_test

import (
	"testing"

	"stashbook-engine/internal/slug"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "simple title", input: "ARK Survival Evolved", want: "ark-survival-evolved"},
		{name: "drops punctuation", input: "Tom & Jerry", want: "tom-jerry"},
		{name: "trims whitespace", input: "  Free Game  ", want: "free-game"},
		// Underscores are dropped by the first pass, before the
		// whitespace pass runs, so they never become hyphens.
		{name: "underscores are dropped", input: "battle_ram_dx", want: "battleramdx"},
		{name: "collapses repeated hyphens", input: "a  --  b", want: "a-b"},
		{name: "trims hyphens", input: "--half-life--", want: "half-life"},
		{name: "keeps digits", input: "Doom 3", want: "doom-3"},
		{name: "drops non-ascii letters", input: "Pokémon", want: "pokmon"},
		{name: "empty input", input: "", want: ""},
		{name: "only punctuation", input: "!!!", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slug.Slugify(tt.input); got != tt.want {
				t.Fatalf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
