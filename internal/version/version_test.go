package version

import "testing"

func TestString(t *testing.T) {
	previousVersion := Version
	previousRevision := Revision
	previousBuildTime := BuildTime
	t.Cleanup(func() {
		Version = previousVersion
		Revision = previousRevision
		BuildTime = previousBuildTime
	})

	Version = "1.2.3"
	Revision = "abc123"
	BuildTime = "2026-09-01T00:00:00Z"

	got := String()
	want := "memovee 1.2.3 (revision abc123, built 2026-09-01T00:00:00Z)"

	if got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
