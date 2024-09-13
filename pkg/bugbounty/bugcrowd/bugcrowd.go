package bugcrowd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/root4loot/goutils/domainutil"
	"github.com/root4loot/goutils/log"
	"github.com/root4loot/goutils/sliceutil"
	"github.com/root4loot/rescope2/pkg/common"
)

type Bugcrowd struct {
	Result common.Result `json:"Result"`
	Auth   string        // _bugcrowd_session=
}

type Target struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URI         string `json:"uri"`
	Category    string `json:"category"`
	InScope     bool   `json:"inScope"`
	SortOrder   int    `json:"sortOrder"`
	Description string `json:"description"`
}

type Scope struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	InScope bool     `json:"inScope"`
	Targets []Target `json:"targets"`
}

type Data struct {
	Scopes []Scope `json:"scope"`
}

type JsonResponse struct {
	Data Data `json:"data"`
}

func (i *Bugcrowd) Run(programURL string, client *http.Client) (*common.Result, error) {
	parsedURL, err := i.ParseURL(programURL)
	if err != nil {
		return nil, err
	}

	i.Result.ProgramDetails = *parsedURL

	if client == nil {
		client = &http.Client{}
	}

	req, err := http.NewRequest("GET", programURL, nil)
	if err != nil {
		return nil, err
	}

	if i.Auth != "" {
		req.Header.Set("Cookie", `_bugcrowd_session="`+i.Auth+`"`)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respB, _ := io.ReadAll(resp.Body)

	log.Debugf("Bugcrowd: Received response with status code %d and body: %s", resp.StatusCode, string(respB))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	re := regexp.MustCompile(`\/changelog\/(\w+-\w+-\w+-\w+-\w+)`)
	UUID := re.FindString(string(respB))
	newURL := "https://bugcrowd.com/engagements/" + parsedURL.ProgramName + UUID + ".json"

	req.URL, err = url.Parse(newURL)
	if err != nil {
		return nil, err
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respB, _ = io.ReadAll(resp.Body)

	var response JsonResponse
	err = json.Unmarshal([]byte(respB), &response)
	if err != nil {
		return nil, err
	}

	for _, scope := range response.Data.Scopes {
		for _, target := range scope.Targets {
			var targetEntry string
			if domainutil.IsDomainName(target.Name) {
				targetEntry = target.Name
			} else {
				targetEntry = target.URI
			}

			if targetEntry == "" {
				continue
			}

			if scope.InScope {
				i.Result.InScope = sliceutil.AppendUnique(i.Result.InScope, targetEntry)
			} else {
				i.Result.OutScope = sliceutil.AppendUnique(i.Result.OutScope, targetEntry)
			}
		}
	}
	return &i.Result, nil
}

func (b *Bugcrowd) ParseURL(rawURL string) (*common.BugBountyProgram, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	if parsedURL.Hostname() != "bugcrowd.com" {
		return nil, fmt.Errorf("invalid domain: %s", parsedURL.Hostname())
	}

	if parsedURL.Path == "" || parsedURL.Path == "/" {
		return nil, fmt.Errorf("invalid program path in URL: %s", parsedURL.Path)
	}

	if !strings.HasPrefix(parsedURL.Path, "/engagements") {
		parsedURL.Path = "/engagements" + parsedURL.Path
	}

	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")

	if len(pathParts) < 2 || pathParts[1] == "" {
		return nil, fmt.Errorf("invalid program path in URL: %s", parsedURL.Path)
	}

	programName := pathParts[1]

	program := &common.BugBountyProgram{
		InputURL:    parsedURL.String(),
		Platform:    "Bugcrowd",
		ProgramName: programName,
		Business:    programName,
		PolicyURL:   "https://" + parsedURL.Hostname() + "/engagements/" + programName,
	}

	return program, nil
}

func (i *Bugcrowd) Serialize() (string, error) {
	jsonData, err := json.Marshal(i.Result)
	if err != nil {
		return "", fmt.Errorf("failed to serialize Result: %w", err)
	}
	return string(jsonData), nil
}
