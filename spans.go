package spans

import ()

import ()

type Range struct {
	Start, End int
}

func (r Range) Len() int {
	return r.End - r.Start
}

type Span interface {
	Start() int
	End() int
}

type SpanWithKind interface {
	Span
	Kind() string
}

func ShortestWithAllKinds(spans []SpanWithKind) (shortest Range) {
	byKind := make(map[string][]SpanWithKind)
	for _, span := range spans {
		byKind[span.Kind()] = append(byKind[span.Kind()], span)
	}
	var kinds []string
	for kind, _ := range byKind {
		kinds = append(kinds, kind)
	}

	combinations := 1
	for _, spans := range byKind {
		combinations *= len(spans)
	}
	for i := 0; i < combinations; i++ {
		try := make([]Span, len(kinds))
		ii := i
		for j, kind := range kinds {
			try[j] = byKind[kind][ii%len(byKind[kind])]
			ii /= len(byKind[kind])
		}
		r := Shortest(try)
		if i == 0 || r.Len() < shortest.Len() {
			shortest = r
		}
	}
	return
}

func Shortest(spans []Span) Range {
	var r Range
	for i, span := range spans {
		if i == 0 || span.Start() < r.Start {
			r.Start = span.Start()
		}
		if i == 0 || span.End() > r.End {
			r.End = span.End()
		}
	}
	return r
}
