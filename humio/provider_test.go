package humio

import (
	"testing"
)

func TestProviderInternalValidation(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
