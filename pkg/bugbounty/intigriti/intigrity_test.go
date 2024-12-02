package intigriti

import (
	"net/http"
	"testing"

	"github.com/root4loot/rescope/pkg/common"
)

// create static intigrity struct

var platform = Intigriti{}

func TestRun(t *testing.T) {
	url := "https://intigriti.com/sqills/sqillscorporatewebsite"
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
		if scope == "*.sqills.com" {
			foundInScope = true
			break
		}
	}

	foundOutScope := false
	for _, scope := range platform.Result.OutScope {
		if scope == "booking.*.sqills.com" {
			foundOutScope = true
			break
		}
	}

	if !foundInScope {
		t.Fatalf("expected *.sqills.com in scope, got %v", platform.Result.InScope)
	}

	if !foundOutScope {
		t.Fatalf("expected booking.*.sqills.com out of scope, got %v", platform.Result.OutScope)
	}
}

func TestParseURL(t *testing.T) {
	tests := []struct {
		inputURL      string
		expectedError bool
		expectedURL   *common.BugBountyProgram
	}{
		{
			inputURL:      "https://app.intigriti.com/programs/sqills/sqillscorporatewebsite/detail",
			expectedError: false,
			expectedURL: &common.BugBountyProgram{
				Platform:    "Intigriti",
				Business:    "sqills",
				ProgramName: "sqillscorporatewebsite",
			},
		},
		{
			inputURL:      "https://intigriti.com/programs/sqills/sqillscorporatewebsite",
			expectedError: false,
			expectedURL: &common.BugBountyProgram{
				Platform:    "Intigriti",
				Business:    "sqills",
				ProgramName: "sqillscorporatewebsite",
			},
		},
		{
			inputURL:      "https://intigriti.com/sqills/sqillscorporatewebsite",
			expectedError: false,
			expectedURL: &common.BugBountyProgram{
				Platform:    "Intigriti",
				Business:    "sqills",
				ProgramName: "sqillscorporatewebsite",
			},
		},
		{
			inputURL:      "https://invalidsite.com/sqills/sqillscorporatewebsite",
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

			// Compare only the relevant fields
			if parsedURL.Platform != test.expectedURL.Platform ||
				parsedURL.Business != test.expectedURL.Business ||
				parsedURL.ProgramName != test.expectedURL.ProgramName {
				t.Fatalf("expected parsed URL Platform: %v, Business: %v, ProgramName: %v, but got Platform: %v, Business: %v, ProgramName: %v",
					test.expectedURL.Platform, test.expectedURL.Business, test.expectedURL.ProgramName,
					parsedURL.Platform, parsedURL.Business, parsedURL.ProgramName)
			}
		}
	}
}
