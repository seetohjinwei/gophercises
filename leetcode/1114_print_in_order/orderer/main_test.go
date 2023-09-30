package orderer

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestSameOrder(t *testing.T) {
	expected := []int{1, 2, 3}

	tests := [][3]int{
		{1, 2, 3},
		{1, 3, 2},
		{2, 1, 3},
		{2, 3, 1},
		{3, 1, 2},
		{3, 2, 1},
	}

	for _, test := range tests {
		orderer := New()
		SameOrder(orderer, test)
		if !slices.Equal(orderer.order, expected) {
			t.Fatalf("Expected %v, got %v", expected, orderer.order)
		}
	}
}
