package bugcrowd

import (
	"net/http"
	"testing"

	"github.com/root4loot/rescope/pkg/common"
)

var platform = Bugcrowd{}

func TestRun(t *testing.T) {
	url := "https://bugcrowd.com/bugcrowd"
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
		if scope == "bugcrowd.com" {
			foundInScope = true
			break
		}
	}

	foundOutScope := false
	for _, scope := range platform.Result.OutScope {
		if scope == "blog.bugcrowd.com" {
			foundOutScope = true
			break
		}
	}

	if !foundInScope {
		t.Fatalf("expected bugcrowd.com in scope, got %v", platform.Result.InScope)
	}

	if !foundOutScope {
		t.Fatalf("expected blog.bugcrowd.com out of scope, got %v", platform.Result.OutScope)
	}
}

func TestParseURL(t *testing.T) {
	tests := []struct {
		inputURL      string
		expectedError bool
		expectedURL   *common.BugBountyProgram
	}{
		{
			inputURL:      "https://bugcrowd.com/engagements/tesla",
			expectedError: false,
			expectedURL: &common.BugBountyProgram{
				Platform:    "Bugcrowd",
				ProgramName: "tesla",
			},
		},
		{
			inputURL:      "https://bugcrowd.com/engagements/unknownprogram",
			expectedError: false,
			expectedURL: &common.BugBountyProgram{
				Platform:    "Bugcrowd",
				ProgramName: "unknownprogram",
			},
		},
		{
			inputURL:      "https://otherdomain.com/engagements/tesla",
			expectedError: true,
			expectedURL:   nil,
		},
		{
			inputURL:      "https://bugcrowd.com/",
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
