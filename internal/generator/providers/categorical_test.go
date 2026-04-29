package providers

import (
	"math/rand/v2"
	"testing"

	"github.com/aapldev00/synthetic_gen_gaas/pkg/genproto"
)

// TestCategoricalFactory verifies that the generator only selects values
// from the provided set and handles whitespace correctly.
func TestCategoricalFactory(t *testing.T) {
	r := rand.New(rand.NewPCG(42, 100))
	rawInput := "RED, GREEN, BLUE"
	expectedSet := map[string]bool{"RED": true, "GREEN": true, "BLUE": true}

	options := []*genproto.Option{
		{Key: "values", Value: &genproto.Option_StringVal{StringVal: rawInput}},
	}

	gen, err := CategoricalFactory(r, options)
	if err != nil {
		t.Fatalf("factory initialization failed: %v", err)
	}

	for i := 0; i < 100; i++ {
		val := gen()

		strVal, ok := val.(string)
		if !ok {
			t.Fatalf("type mismatch: expected string, got %T", val)
		}

		if !expectedSet[strVal] {
			t.Errorf("value %s not found in original set %v", strVal, expectedSet)
		}
	}
}
