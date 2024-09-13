package hackenproof

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/root4loot/goutils/log"
	"github.com/root4loot/goutils/sliceutil"
	"github.com/root4loot/rescope2/pkg/common"
)

type HackenProof struct {
	Result common.Result `json:"Result"`
	Auth   string        // cookie: _hackenproof_session
}

type ScopeDetails struct {
	Type              string `json:"type"`
	Target            string `json:"target"`
	TargetDescription string `json:"target_description"`
	Severity          string `json:"severity"`
	RewardType        string `json:"reward_type"`
	OutOfScope        bool   `json:"out_of_scope"`
}

type ProgramData struct {
	Scopes []ScopeDetails `json:"scopes"`
}

func (i *HackenProof) Run(programURL string, client *http.Client) (*common.Result, error) {
	parsedURL, err := i.ParseURL(programURL)
	if err != nil {
		return nil, err
	}

	i.Result.ProgramDetails = *parsedURL

	req, err := http.NewRequest("GET", parsedURL.InputURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Cookie", "_hackenproof_session="+i.Auth)

	if client == nil {
		client = &http.Client{}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Debugf("HackenProof: Received response with status code %d and body: %s", resp.StatusCode, string(body))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad response status: %s", resp.Status)
	}

	var programData ProgramData
	err = json.Unmarshal(body, &programData)
	if err != nil {
		return nil, err
	}

	for _, scope := range programData.Scopes {
		if strings.ToLower(scope.Target) == "" {
			continue
		}
		if scope.OutOfScope {
			i.Result.OutScope = sliceutil.AppendUnique(i.Result.OutScope, strings.ToLower(scope.Target))
		} else {
			i.Result.InScope = sliceutil.AppendUnique(i.Result.InScope, strings.ToLower(scope.Target))
		}
	}

	return &i.Result, nil
}

func (b *HackenProof) ParseURL(rawURL string) (*common.BugBountyProgram, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	if parsedURL.Hostname() != "hackenproof.com" {
		return nil, fmt.Errorf("invalid domain: %s", parsedURL.Hostname())
	}

	if !strings.HasPrefix(parsedURL.Path, "/bug-bounty-programs-list") {
		parsedURL.Path = "/bug-bounty-programs-list" + parsedURL.Path
	}

	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")

	if len(pathParts) != 2 || pathParts[1] == "" {
		return nil, fmt.Errorf("invalid program path in URL: %s", parsedURL.Path)
	}

	programName := pathParts[1]

	program := &common.BugBountyProgram{
		InputURL:    parsedURL.String(),
		Platform:    "HackenProof",
		ProgramName: programName,
		Business:    programName,
		PolicyURL:   "https://" + parsedURL.Hostname() + "bug-bounty-programs-list/" + programName,
	}

	return program, nil
}

func (i *HackenProof) Serialize() (string, error) {
	jsonData, err := json.Marshal(i.Result)
	if err != nil {
		return "", fmt.Errorf("failed to serialize Result: %w", err)
	}
	return string(jsonData), nil
}
