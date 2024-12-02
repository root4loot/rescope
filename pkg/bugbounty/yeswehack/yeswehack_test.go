package yeswehack

import (
	"net/http"
	"testing"

	"github.com/root4loot/rescope/pkg/common"
)

var platform = YesWeHack{}

func TestRun(t *testing.T) {
	url := "https://yeswehack.com/programs/legapass-bug-bounty-program"
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
		if scope == "https://bounty.legapass.com" {
			foundInScope = true
			break
		}
	}

	foundOutScope := false
	for _, scope := range platform.Result.OutScope {
		if scope == "app.legapass.com" {
			foundOutScope = true
			break
		}
	}

	if !foundInScope {
		t.Fatalf("expected hackerone.com in scope, got %v", platform.Result.InScope)
	}

	if !foundOutScope {
		t.Fatalf("expected app.legapass.com out of scope, got %v", platform.Result.OutScope)
	}
}

func TestParseURL(t *testing.T) {
	tests := []struct {
		inputURL        string
		expectedError   bool
		expectedProgram *common.BugBountyProgram
	}{
		{
			inputURL:      "https://yeswehack.com/programs/legapass-bug-bounty-program",
			expectedError: false,
			expectedProgram: &common.BugBountyProgram{
				Platform:    "YesWeHack",
				ProgramName: "legapass-bug-bounty-program",
			},
		},
		{
			inputURL:        "https://otherdomain.com/programs/legapass-bug-bounty-program",
			expectedError:   true,
			expectedProgram: nil,
		},
		{
			inputURL:        "https://yeswehack.com/otherpath/legapass-bug-bounty-program",
			expectedError:   true,
			expectedProgram: nil,
		},
		{
			inputURL:        "https://yeswehack.com/programs/",
			expectedError:   true,
			expectedProgram: nil,
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

			if parsedURL.Platform != test.expectedProgram.Platform ||
				parsedURL.ProgramName != test.expectedProgram.ProgramName {
				t.Fatalf("expected parsed URL Platform: %v, ProgramName: %v, but got Platform: %v, ProgramName: %v",
					test.expectedProgram.Platform, test.expectedProgram.ProgramName,
					parsedURL.Platform, parsedURL.ProgramName)
			}
		}
	}
}
