package providers

import (
	"math/rand/v2"
	"testing"

	"github.com/google/uuid"
)

// TestUUIDFactory verifies that the UUID provider produces valid UUID v4 strings
// and maintains the expected functional signature.
func TestUUIDFactory(t *testing.T) {
	// 1. Setup: Entropy source (ignored by the provider but required by the signature)
	r := rand.New(rand.NewPCG(1, 2))

	// 2. Call the factory (Planning Phase)
	gen, err := UUIDFactory(r, nil)
	if err != nil {
		t.Fatalf("UUIDFactory returned an unexpected error: %v", err)
	}

	// 3. Validation Loop (Execution Phase)
	// We generate multiple IDs to ensure consistency and validity.
	for i := 0; i < 50; i++ {
		val := gen()

		// Verify semantic type preservation (must be string)
		strVal, ok := val.(string)
		if !ok {
			t.Fatalf("Iteration %d: expected string type, got %T", i, val)
		}

		// Verify UUID format validity using the official Google UUID library
		_, err := uuid.Parse(strVal)
		if err != nil {
			t.Errorf("Iteration %d: generated an invalid UUID string: %s", i, strVal)
		}
	}
}
