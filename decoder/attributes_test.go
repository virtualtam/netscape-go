package decoder

import "testing"

func assertAttributesEqual(t *testing.T, got map[string]string, want map[string]string) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("want %d attributes, got %d", len(want), len(got))
		return
	}

	for attr, wantValue := range want {
		if got[attr] != wantValue {
			t.Errorf("want attribute %q value %q, got %q", attr, wantValue, got[attr])
		}
	}
}
