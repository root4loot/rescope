package rescope

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/root4loot/goutils/domainutil"
	"github.com/root4loot/goutils/log"

	"github.com/root4loot/rescope2/pkg/bugbounty/bugcrowd"
	"github.com/root4loot/rescope2/pkg/bugbounty/hackenproof"
	"github.com/root4loot/rescope2/pkg/bugbounty/hackerone"
	"github.com/root4loot/rescope2/pkg/bugbounty/intigriti"
	"github.com/root4loot/rescope2/pkg/bugbounty/yeswehack"
	"github.com/root4loot/rescope2/pkg/common"
)

type Result interface {
	Serialize() (string, error)
}

type BugBountyProgram interface {
	Run(url string, client *http.Client) (*common.Result, error)
	ParseURL(url string) (*common.BugBountyProgram, error)
	Serialize() (string, error)
}

type Options struct {
	Client          *http.Client
	AuthHackerOne   string
	AuthIntigriti   string
	AuthBugcrowd    string
	AuthHackenProof string
	AuthYesWeHack   string
	Debug           bool
}

func DefaultOptions() *Options {
	return &Options{
		Client:          &http.Client{},
		AuthHackerOne:   "",
		AuthIntigriti:   "",
		AuthBugcrowd:    "",
		AuthHackenProof: "",
		AuthYesWeHack:   "",
		Debug:           false,
	}
}

// Run function with validation

func Run(url string, options *Options) (*common.Result, error) {

	if options.Debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Debug("Running rescope", "target", url, "options", options)

	platform, err := IdentifyPlatform(url, options)
	if err != nil {
		return nil, errors.Wrap(err, "unsupported or invalid URL")
	}

	result, err := platform.Run(url, options.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run platform")
	}

	result.ProgramDetails.FetchedAt = time.Now().Format(time.RFC3339)
	return result, nil
}

func IsBugBountyURL(bugbountyURL string) bool {
	u, err := url.Parse(bugbountyURL)
	if err != nil {
		return false
	}

	rootDomain := domainutil.GetRootDomain(u.Hostname())

	switch rootDomain {
	case "intigriti.com", "hackerone.com", "yeswehack.com", "bugcrowd.com", "hackenproof.com":
		return true
	default:
		return false
	}
}

func IdentifyPlatform(bugbountyURL string, options *Options) (BugBountyProgram, error) {
	u, err := url.Parse(bugbountyURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse URL")
	}

	rootDomain := domainutil.GetRootDomain(u.Hostname())

	switch rootDomain {
	case "intigriti.com":
		return &intigriti.Intigriti{Auth: options.AuthIntigriti}, nil
	case "hackerone.com":
		return &hackerone.HackerOne{Auth: options.AuthHackerOne}, nil
	case "yeswehack.com":
		return &yeswehack.YesWeHack{Auth: options.AuthYesWeHack}, nil
	case "bugcrowd.com":
		return &bugcrowd.Bugcrowd{Auth: options.AuthBugcrowd}, nil
	case "hackenproof.com":
		return &hackenproof.HackenProof{Auth: options.AuthHackenProof}, nil
	default:
		return nil, fmt.Errorf("unsupported bug bounty platform for URL: %s", bugbountyURL)
	}
}
