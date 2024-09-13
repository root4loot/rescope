package hackerone

import (
	"bytes"
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

type HackerOne struct {
	Result common.Result `json:"Result"`
	Auth   string        // authorization bearer token
}

func (i *HackerOne) Run(programURL string, client *http.Client) (*common.Result, error) {
	if client == nil {
		client = &http.Client{}
	}

	parsedURL, err := i.ParseURL(programURL)
	if err != nil {
		return nil, err
	}

	i.Result.ProgramDetails = *parsedURL

	var data = []byte(`{
		"query":"query Team_assets($first_0:Int!) {query {id,...F0}} fragment F0 on Query {_teamAgUhl:team(handle:\"` + parsedURL.ProgramName + `\") {handle,_structured_scope_versions2ZWKHQ:structured_scope_versions(archived:false) {max_updated_at},_structured_scopeszxYtW:structured_scopes(first:$first_0,archived:false,eligible_for_submission:true) {edges {node {asset_type, asset_identifier}},pageInfo {hasNextPage,hasPreviousPage}},_structured_scopes3FF98f:structured_scopes(first:$first_0,archived:false,eligible_for_submission:false) {edges {node {asset_type,asset_identifier,},},},},}",
		"variables":{
		   "first_0":1337
		}
	 }`)

	req, err := http.NewRequest("GET", programURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	resB, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Debugf("HackerOne: Received response with status code %d and body: %s", resp.StatusCode, string(resB))

	if resp.StatusCode != 200 {
		return nil, err
	}

	hostsession, csrf, err := getSessionAndCSRF(*client)
	if err != nil {
		return nil, err
	}

	req, _ = http.NewRequest("POST", "https://hackerone.com/graphql?", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")
	req.Header.Set("Cookie", hostsession)
	req.Header.Set("X-Csrf-Token", csrf)

	if i.Auth != "" {
		req.Header.Set("X-Auth-Token", i.Auth)
	} else {
		log.Debug("hackerone: No token provided, running unauthenticated...")
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resB, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`\"edges":\[(.*?)\]`)
	scopeSplit := re.FindAllString(string(resB), -1)
	re = regexp.MustCompile(`asset_type":"(URL|CIDR|IP|IP-RANGE|RANGE)","asset_identifier":"(.*?)"`)

	for _, match := range re.FindAllStringSubmatch(scopeSplit[0], -1) {
		if strings.ToLower(match[2]) == "" {
			continue
		}
		i.Result.InScope = append(i.Result.InScope, match[2])
	}

	for _, match := range re.FindAllStringSubmatch(scopeSplit[1], -1) {
		if strings.ToLower(match[2]) == "" {
			continue
		}
		i.Result.OutScope = append(i.Result.OutScope, match[2])
	}
	return &i.Result, err
}

func (r *HackerOne) ParseURL(rawURL string) (*common.BugBountyProgram, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	if parsedURL.Hostname() != "hackerone.com" {
		return nil, fmt.Errorf("invalid domain: %s", parsedURL.Hostname())
	}

	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathParts) < 1 || pathParts[0] == "" {
		return nil, fmt.Errorf("invalid program path in URL: %s", parsedURL.Path)
	}

	program := pathParts[0]

	urlStruct := &common.BugBountyProgram{
		InputURL:    rawURL,
		Platform:    "HackerOne",
		ProgramName: program,
		Business:    program,
		PolicyURL:   "https://" + parsedURL.Host + "/" + program,
	}

	return urlStruct, nil
}

func (h *HackerOne) Serialize() (string, error) {
	jsonData, err := json.Marshal(h.Result)
	if err != nil {
		return "", fmt.Errorf("failed to serialize Result: %w", err)
	}
	return string(jsonData), nil
}

func getSessionAndCSRF(client http.Client) (hostsession, csrfToken string, err error) {
	req, err := http.NewRequest("GET", "https://hackerone.com/security", nil)
	if err != nil {
		return hostsession, csrfToken, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return hostsession, csrfToken, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return hostsession, csrfToken, err
	}

	cookies := resp.Header["Set-Cookie"]
	for _, cookie := range cookies {
		if strings.HasPrefix(cookie, "__Host-session") {
			hostsession = strings.Split(cookie, ";")[0]
		}
	}

	r := regexp.MustCompile(`<meta name="csrf-token" content="([\w+\/=]+)`)
	m := r.FindStringSubmatch(string(body))
	csrfToken = m[1]

	return hostsession, csrfToken, err
}
