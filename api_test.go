package coven

import (
	"testing"
)

func TestCache(t *testing.T) {
	type foo struct {
		A int
	}

	type bar struct {
		A int
	}

	_ = NewConverter(new(foo), new(bar))
	_ = NewConverter(new(foo), new(bar))

	if len(createdConverters) != 1 {
		t.Fatalf("cache fail")
	}
}
