package command

import "testing"

func TestExitCodes(t *testing.T) {
	tests := []struct {
		category Category
		want     int
	}{
		{CategoryUsage, ExitUsage},
		{CategoryPrerequisite, ExitPrerequisite},
		{CategoryContract, ExitContract},
		{CategoryConfiguration, ExitConfiguration},
		{CategoryOwnership, ExitOwnership},
		{CategorySecret, ExitSecret},
		{CategoryProcess, ExitProcess},
		{CategoryNetwork, ExitNetwork},
		{CategoryVerification, ExitVerification},
		{CategoryActivation, ExitActivation},
		{CategoryRollback, ExitRollback},
		{CategoryInternal, ExitInternal},
	}

	for _, test := range tests {
		t.Run(string(test.category), func(t *testing.T) {
			if got := ExitCode(test.category); got != test.want {
				t.Fatalf("ExitCode(%q) = %d, want %d", test.category, got, test.want)
			}
		})
	}
}

func TestUnknownCategoryUsesInternalExitCode(t *testing.T) {
	if got := ExitCode(Category("unknown")); got != ExitInternal {
		t.Fatalf("ExitCode(unknown) = %d, want %d", got, ExitInternal)
	}
}
