package pipeline

import "testing"

func TestDecode(t *testing.T) {
	t.Run("valid json", func(t *testing.T) {
		raw := []byte(`{"id":"CVE-2021-44228","package":"log4j","severity":"critical","cvss":10}`)

		got, err := Decode(raw)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := Finding{ID: "CVE-2021-44228", Package: "log4j", Severity: Critical, CVSS: 10}
		if got != want {
			t.Errorf("got %+v, want %+v", got, want)
		}

	})

	t.Run("malformed json returns an error", func(t *testing.T) {
		_, err := Decode([]byte(`{not json`))
		if err == nil {
			t.Fatal("expected an error for malformed json, got nil")
		}
	})
}
