package intigriti

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

type Intigriti struct {
	Result common.Result `json:"Result"`
	Auth   string        // authorization bearer token
}

func (i *Intigriti) Run(programURL string, client *http.Client) (*common.Result, error) {
	parsedURL, err := i.ParseURL(programURL)
	if err != nil {
		return nil, err
	}

	i.Result.ProgramDetails = *parsedURL

	if client == nil {
		client = &http.Client{}
	}

	tryFetchPublicScope := true

	if i.Auth != "" {
		log.Debug("Token provided, attempting to fetch private scope data")
		privateScopeDetails, err := fetchPrivateScope(*parsedURL, i.Auth, *client)
		if err == nil && privateScopeDetails != nil {
			processPrivateScope(&i.Result, privateScopeDetails)
			tryFetchPublicScope = false
		} else {
			log.Warn("Failed to fetch private scope data, will attempt public scope")
		}
	}

	if tryFetchPublicScope {
		log.Debug("Fetching public scope data")
		publicProgramDetail, err := fetchPublicScope(*parsedURL, client)
		if err != nil {
			return nil, err
		}
		processPublicScope(&i.Result, publicProgramDetail)
	}

	return &i.Result, nil
}

func (i *Intigriti) ParseURL(rawURL string) (*common.BugBountyProgram, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if u.Host != "intigriti.com" && u.Host != "app.intigriti.com" {
		return nil, fmt.Errorf("invalid domain '%s': %w", u.Hostname(), err)
	}

	if u.Host == "app.intigriti.com" {
		u.Host = "intigriti.com"
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid URL path '%s': %w", u.Path, err)
	}

	if parts[0] != "programs" {
		parts = append([]string{"programs"}, parts...)
	}

	business := parts[1]
	program := parts[2]

	return &common.BugBountyProgram{
		InputURL:    u.String(),
		Platform:    "Intigriti",
		Business:    business,
		ProgramName: program,
		PolicyURL:   "https://" + u.Hostname() + "/" + business + "/" + program + "/detail",
	}, nil
}

func (i *Intigriti) Serialize() (string, error) {
	jsonData, err := json.Marshal(i.Result)
	if err != nil {
		return "", fmt.Errorf("failed to serialize Result: %w", err)
	}
	return string(jsonData), nil
}

func fetchPrivateProgramList(token string, client *http.Client) (*PrivateProgramList, error) {
	endpoint := "https://api.intigriti.com/external/researcher/v1/programs?following=false"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Debugf("Intigriti: Received response with status code %d and body: %s", resp.StatusCode, string(respB))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	var privateProgramList PrivateProgramList

	err = json.Unmarshal(respB, &privateProgramList)
	if err != nil {
		return nil, err
	}

	return &privateProgramList, nil
}

func fetchPrivateScope(url common.BugBountyProgram, token string, client http.Client) (*PrivateProgramDetail, error) {
	privateProgramList, err := fetchPrivateProgramList(token, &client)
	if err != nil {
		return nil, err
	}

	for _, program := range privateProgramList.Records {
		if program.Handle == url.ProgramName {
			url.ProgramName = program.ID
			endpoint := fmt.Sprintf("https://api.intigriti.com/external/researcher/v1/programs/%s/", program.ID)

			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				return nil, err
			}

			req.Header.Add("Accept", "application/json")
			req.Header.Add("Authorization", "Bearer "+token)

			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}

			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("expected status code 200, got %d", resp.StatusCode)
			}

			respB, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			var privateProgramDetail PrivateProgramDetail

			err = json.Unmarshal(respB, &privateProgramDetail)
			if err != nil {
				return nil, err
			}

			return &privateProgramDetail, nil
		}
	}

	return nil, nil
}

func fetchPublicScope(program common.BugBountyProgram, client *http.Client) (*PublicProgramDetail, error) {
	endpoint := fmt.Sprintf("https://app.intigriti.com/api/core/public/programs/%s/%s", program.Business, program.ProgramName)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var publicProgramDetail PublicProgramDetail

	err = json.Unmarshal(respB, &publicProgramDetail)
	if err != nil {
		return nil, err
	}

	return &publicProgramDetail, nil
}

func processPrivateScope(Result *common.Result, scopeDetails *PrivateProgramDetail) {
	for _, content := range scopeDetails.Domains.Content {
		if content.Endpoint == "" {
			continue
		}
		if content.Type.Value == "Url" {
			Result.InScope = sliceutil.AppendUnique(Result.InScope, content.Endpoint)
		}
	}
}

func processPublicScope(Result *common.Result, publicProgramDetail *PublicProgramDetail) {
	for _, domain := range publicProgramDetail.Domains {
		for _, content := range domain.Content {
			if content.Endpoint == "" {
				continue
			}
			if content.BountyTierID != 5 {
				Result.InScope = sliceutil.AppendUnique(Result.InScope, content.Endpoint)
			} else {
				Result.OutScope = sliceutil.AppendUnique(Result.OutScope, content.Endpoint)
			}
		}
	}
}
