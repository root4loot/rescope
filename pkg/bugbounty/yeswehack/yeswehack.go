package yeswehack

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/root4loot/goutils/log"
	"github.com/root4loot/rescope2/pkg/common"
)

type YesWeHack struct {
	Result common.Result `json:"Result"`
	Auth   string        // authorization bearer token
}

func (i *YesWeHack) Run(programURL string, client *http.Client) (*common.Result, error) {
	parsedURL, err := i.ParseURL(programURL)
	if err != nil {
		return nil, err
	}

	i.Result.ProgramDetails = *parsedURL

	req, err := http.NewRequest("GET", "https://api.yeswehack.com/programs/"+parsedURL.ProgramName, nil)
	if err != nil {
		return nil, err
	}

	if i.Auth != "" {
		req.Header.Set("Authorization", "Bearer "+i.Auth)
	}

	if client == nil {
		client = &http.Client{}
	}

	req.URL, err = url.Parse(programURL)
	if err != nil {
		return nil, err
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

	log.Debugf("YesWeHack: Received response with status code %d and body: %s", resp.StatusCode, string(body))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	inScopeRegex := regexp.MustCompile(`"scope":"([^"]*?(\*\.[a-zA-Z0-9.-]+|https?:\/\/[a-zA-Z0-9.-]+|[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})(?:\\[^"]*)?)"`)
	inScopeMatches := inScopeRegex.FindAllStringSubmatch(string(body), -1)

	for _, match := range inScopeMatches {
		if len(match) > 1 {
			cleanedURL := strings.ReplaceAll(strings.TrimSpace(match[1]), `\`, ``)
			if cleanedURL != "" {
				i.Result.InScope = append(i.Result.InScope, cleanedURL)
			}
		}
	}

	outScopeRegex := regexp.MustCompile(`"out_of_scope":\[(.*?)\],`)
	outScopeTextMatches := outScopeRegex.FindStringSubmatch(string(body))

	if len(outScopeTextMatches) > 1 {
		cleanedOutScopeText := strings.ReplaceAll(outScopeTextMatches[1], `\"`, `"`)
		urlRegex := regexp.MustCompile(`(?:\*?\.)?[a-zA-Z0-9-]+(?:\.[a-zA-Z]{2,})(?:\.[a-zA-Z]{2,})?(?:$begin:math:text$[a-z|]+$end:math:text$)?`)
		outScopeURLs := urlRegex.FindAllString(cleanedOutScopeText, -1)

		for _, url := range outScopeURLs {
			if url != "" {
				i.Result.OutScope = append(i.Result.OutScope, url)
			}
		}
	}

	return &i.Result, nil
}

func (i *YesWeHack) ParseURL(rawURL string) (*common.BugBountyProgram, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ProgramName: %w", err)
	}

	if parsedURL.Hostname() != "yeswehack.com" {
		return nil, fmt.Errorf("invalid domain: %s", parsedURL.Hostname())
	}

	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathParts) != 2 || pathParts[0] != "programs" {
		return nil, fmt.Errorf("invalid program path in ProgramName: %s", parsedURL.Path)
	}

	program := pathParts[1]

	programStruct := &common.BugBountyProgram{
		InputURL:    rawURL,
		Platform:    "YesWeHack",
		ProgramName: program,
		Business:    program,
		PolicyURL:   "https://" + parsedURL.Hostname() + "/programs/" + program,
	}

	return programStruct, nil
}

func (y *YesWeHack) Serialize() (string, error) {
	jsonData, err := json.Marshal(y.Result)
	if err != nil {
		return "", fmt.Errorf("failed to serialize Result: %w", err)
	}
	return string(jsonData), nil
}
