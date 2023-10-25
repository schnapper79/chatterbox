package types

import "testing"

func Test_NewModelRequestWithDefaults(t *testing.T) {
	mr := NewModelRequestWithDefaults()
	if mr.ContextSize != 4096 {
		t.Errorf("ContextSize should be 4096, got %d", mr.ContextSize)
	}
}

func Test_NewPredictionRequestWithDefaults(t *testing.T) {
	pr := NewPredictionRequestWithDefaults()
	if pr.TopP != 0.95 {
		t.Errorf("TopP should be 0.95, got %v", pr.TopP)
	}
}
