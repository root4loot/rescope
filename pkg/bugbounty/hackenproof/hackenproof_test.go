package hackenproof

import (
	"net/http"
	"testing"

	"github.com/root4loot/rescope2/pkg/common"
)

var platform = HackenProof{}

func TestRun(t *testing.T) {
	url := "https://hackenproof.com/bug-bounty-programs-list/internet-computer-protocol"
	client := &http.Client{}

	Result, err := platform.Run(url, client)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if Result == nil {
		t.Fatalf("expected a valid Result, got nil")
	}

	foundInScope := false
	for _, scope := range platform.Result.InScope {
		if scope == "boundary.ic0.app" {
			foundInScope = true
			break
		}
	}

	if !foundInScope {
		t.Fatalf("expected boundary.ic0.app in scope, got %v", platform.Result.InScope)
	}
}

func TestParseURL(t *testing.T) {
	tests := []struct {
		inputURL      string
		expectedError bool
		expectedURL   *common.BugBountyProgram
	}{
		{
			inputURL:      "https://hackenproof.com/internet-computer-protocol",
			expectedError: false,
			expectedURL: &common.BugBountyProgram{
				Platform:    "HackenProof",
				Business:    "internet-computer-protocol",
				ProgramName: "internet-computer-protocol",
			},
		},
		{
			inputURL:      "https://hackenproof.com/bug-bounty-programs-list/internet-computer-protocol",
			expectedError: false,
			expectedURL: &common.BugBountyProgram{
				Platform:    "HackenProof",
				Business:    "internet-computer-protocol",
				ProgramName: "internet-computer-protocol",
			},
		},
		{
			inputURL:      "https://invalidsite.com/internet-computer-protocol",
			expectedError: true,
			expectedURL:   nil,
		},
		{
			inputURL:      "https://hackenproof.com/invalid/url/structure",
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
				parsedURL.Business != test.expectedURL.Business ||
				parsedURL.ProgramName != test.expectedURL.ProgramName {
				t.Fatalf("expected parsed URL %v, but got %v", test.expectedURL, parsedURL)
			}
		}
	}
}
