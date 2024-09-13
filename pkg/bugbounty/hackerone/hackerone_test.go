package hackerone

import (
	"net/http"
	"testing"

	"github.com/root4loot/rescope2/pkg/common"
)

var platform = HackerOne{}

func TestRun(t *testing.T) {
	client := &http.Client{}

	Result, err := platform.Run("https://hackerone.com/security", client)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if Result == nil {
		t.Fatalf("expected a valid Result, got nil")
	}

	foundInScope := false
	for _, scope := range platform.Result.InScope {
		if scope == "hackerone.com" {
			foundInScope = true
			break
		}
	}

	foundOutScope := false
	for _, scope := range platform.Result.OutScope {
		if scope == "support.hackerone.com" {
			foundOutScope = true
			break
		}
	}

	if !foundInScope {
		t.Fatalf("expected hackerone.com in scope, got %v", platform.Result.InScope)
	}

	if !foundOutScope {
		t.Fatalf("expected support.hackerone.com out of scope, got %v", platform.Result.OutScope)
	}
}

func TestParseURL(t *testing.T) {
	tests := []struct {
		inputURL      string
		expectedError bool
		expectedURL   *common.BugBountyProgram
	}{
		{
			inputURL:      "https://hackerone.com/security",
			expectedError: false,
			expectedURL: &common.BugBountyProgram{
				Platform:    "HackerOne",
				ProgramName: "security",
			},
		},
		{
			inputURL:      "https://hackerone.com/BugBountyProgram",
			expectedError: false,
			expectedURL: &common.BugBountyProgram{
				Platform:    "HackerOne",
				ProgramName: "BugBountyProgram",
			},
		},
		{
			inputURL:      "https://otherdomain.com/security",
			expectedError: true,
			expectedURL:   nil,
		},
		{
			inputURL:      "https://hackerone.com/",
			expectedError: true,
			expectedURL:   nil,
		},
		{
			inputURL:      "https://hackerone.com//",
			expectedError: true,
			expectedURL:   nil,
		},
	}

	for _, test := range tests {
		parsedURL, err := platform.ParseURL(test.inputURL)
		if test.expectedError {
			if err == nil {
				t.Fatalf("expected an error for URL %s, but got none", test.inputURL)
			}
		} else {
			if err != nil {
				t.Fatalf("did not expect an error for URL %s, but got: %v", test.inputURL, err)
			}
			if parsedURL.Platform != test.expectedURL.Platform ||
				parsedURL.ProgramName != test.expectedURL.ProgramName {
				t.Fatalf("expected parsed URL %v, but got %v", test.expectedURL, parsedURL)
			}
		}
	}
}
