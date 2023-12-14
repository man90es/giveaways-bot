package utils

import (
	"testing"
)

const iterations = 1000

func TestRandomChoice(t *testing.T) {
	list := []string{"a", "b", "c", "d", "e"}
	counts := make(map[string]int)

	for i := 0; i < iterations; i++ {
		selected, _ := RandomChoice(list)
		counts[selected]++
	}

	expectedFraction := 1.0 / float64(len(list))
	for _, element := range list {
		actualFraction := float64(counts[element]) / float64(iterations)
		if !approximatelyEqual(actualFraction, expectedFraction, 0.1) {
			t.Errorf("Fraction of %v is not approximately equal. Expected: %v, Actual: %v", element, expectedFraction, actualFraction)
		}
	}

	_, err := RandomChoice([]string{})
	if err == nil {
		t.Errorf("RandomChoice on an empty list should return an error, got: %v", err)
	}
}

func approximatelyEqual(a, b, tolerance float64) bool {
	return (a-b) < tolerance && (b-a) < tolerance
}
