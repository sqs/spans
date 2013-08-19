package spans

import (
	"flag"
	"strings"
	"testing"
)

var q = flag.String("test.q", "", "only run tests whose label contains this string")

type span struct {
	start, end int
}

type spanWithKind struct {
	span
	kind string
}

func (s span) Start() int           { return s.start }
func (s span) End() int             { return s.end }
func (s spanWithKind) Kind() string { return s.kind }

func TestShortest(t *testing.T) {
	tests := map[string]struct {
		spans    []Span
		shortest Range
	}{
		"empty": {},
		"1 span": {
			spans:    []Span{span{0, 1}},
			shortest: Range{0, 1},
		},
		"2 non-overlapping spans": {
			spans:    []Span{span{0, 1}, span{2, 3}},
			shortest: Range{0, 3},
		},
		"2 overlapping spans": {
			spans:    []Span{span{0, 2}, span{1, 3}},
			shortest: Range{0, 3},
		},
	}
	for label, test := range tests {
		shortest := Shortest(test.spans)
		if test.shortest != shortest {
			t.Errorf("%s: want %+v, got %+v", label, test.shortest, shortest)
		}
	}
}

func TestShortestWithAllKinds(t *testing.T) {
	tests := map[string]struct {
		spans    []SpanWithKind
		shortest Range
	}{
		"empty": {},
		"1 span": {
			spans:    []SpanWithKind{spanWithKind{span{0, 1}, "A"}},
			shortest: Range{0, 1},
		},
		"2 non-overlapping spans (same size), 1 kind": {
			spans:    []SpanWithKind{spanWithKind{span{0, 1}, "A"}, spanWithKind{span{2, 3}, "A"}},
			shortest: Range{0, 1},
		},
		"2 non-overlapping spans (diff size), 1 kind": {
			spans:    []SpanWithKind{spanWithKind{span{0, 2}, "A"}, spanWithKind{span{2, 3}, "A"}},
			shortest: Range{2, 3},
		},
		"2 overlapping spans (same size), 1 kind": {
			spans:    []SpanWithKind{spanWithKind{span{0, 2}, "A"}, spanWithKind{span{1, 3}, "A"}},
			shortest: Range{0, 2},
		},
		"2 overlapping spans (diff size), 1 kind": {
			spans:    []SpanWithKind{spanWithKind{span{0, 3}, "A"}, spanWithKind{span{1, 3}, "A"}},
			shortest: Range{1, 3},
		},

		// >1 kinds
		"2 non-overlapping spans (same size), 2 kind": {
			spans:    []SpanWithKind{spanWithKind{span{0, 2}, "A"}, spanWithKind{span{3, 5}, "B"}},
			shortest: Range{0, 5},
		},
		"2 non-overlapping spans (diff size), 2 kind": {
			spans:    []SpanWithKind{spanWithKind{span{0, 1}, "A"}, spanWithKind{span{3, 5}, "B"}},
			shortest: Range{0, 5},
		},
		"2 overlapping spans (same size), 2 kind": {
			spans:    []SpanWithKind{spanWithKind{span{0, 2}, "A"}, spanWithKind{span{1, 3}, "B"}},
			shortest: Range{0, 3},
		},
		"2 overlapping spans (diff size), 2 kind": {
			spans:    []SpanWithKind{spanWithKind{span{0, 3}, "A"}, spanWithKind{span{1, 3}, "B"}},
			shortest: Range{0, 3},
		},

		// >1 kinds, >1 spans per kind
		"multi 1": {
			spans:    []SpanWithKind{spanWithKind{span{0, 3}, "A"}, spanWithKind{span{4, 5}, "A"}, spanWithKind{span{5, 6}, "B"}},
			shortest: Range{4, 6},
		},
		"multi 2": {
			spans:    []SpanWithKind{spanWithKind{span{0, 1}, "A"}, spanWithKind{span{4, 5}, "A"}, spanWithKind{span{5, 6}, "B"}, spanWithKind{span{0, 1}, "B"}},
			shortest: Range{0, 1},
		},
		"multi 3": {
			spans:    []SpanWithKind{spanWithKind{span{0, 1}, "A"}, spanWithKind{span{3, 4}, "A"}, spanWithKind{span{5, 6}, "A"}, spanWithKind{span{1, 3}, "B"}, spanWithKind{span{5, 7}, "B"}},
			shortest: Range{5, 7},
		},
		"multi 4": {
			spans:    []SpanWithKind{spanWithKind{span{0, 1}, "A"}, spanWithKind{span{1, 2}, "A"}, spanWithKind{span{1, 2}, "B"}, spanWithKind{span{2, 3}, "B"}, spanWithKind{span{3, 4}, "C"}, spanWithKind{span{4, 5}, "C"}},
			shortest: Range{1, 4},
		},
	}
	for label, test := range tests {
		if !strings.Contains(label, *q) {
			continue
		}
		shortest := ShortestWithAllKinds(test.spans)
		if test.shortest != shortest {
			t.Errorf("%s: want %+v, got %+v", label, test.shortest, shortest)
		}
	}
}
