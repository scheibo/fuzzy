package fuzzy

import (
	"testing"
)

func TestMatch(t *testing.T) {
	cs := []string{"foo", "bar", "baz", "hello", "world", "cart", "camp", "code", "fuzz"}
	tests := []struct {
		s, m string
		ts   float64
	}{
		{"fo", "foo", 0.5},
		{"fooo", "foo", 0.5},
		{"fo", "foo", 0.5},
		{"hello", "hello", 1.0},
		{"orld", "world", 0.8},
		{"od", "code", 0.1},
		{"ba", "baz", 0.5},
		{"ba", "baz", 0.2},
	}
	for _, tt := range tests {
		m, ts := Match(tt.s, cs)
		if m != tt.m || ts < tt.ts {
			t.Errorf("Match(%s, %v): got (%s, %.3f), want: (%s, >=%.3f)", tt.s, cs, m, ts, tt.m, tt.ts)
		}
	}
}
