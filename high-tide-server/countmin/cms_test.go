package countmin

import "testing"

func TestCountMinSketch(t *testing.T) {
	// Initialize your sketch
	cms := NewCountMinSketch(3, 4) // example dims

	// Perform an action "n" times
	n := 5
	for range n {
		cms.Update("user_123")
	}

	// Verify the result
	count := cms.PointQuery("user_123")
	if count != n {
		t.Errorf("Expected count %d, got %d", n, count)
	}
}
