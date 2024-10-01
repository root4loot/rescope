package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/root4loot/goutils/fileutil"
	"github.com/root4loot/goutils/iputil"
	"github.com/root4loot/goutils/log"
	"github.com/root4loot/goutils/sliceutil"
	"github.com/root4loot/goutils/urlutil"
	"github.com/root4loot/rescope2/config"
	"github.com/root4loot/rescope2/pkg/common"
	"github.com/root4loot/rescope2/pkg/rescope"
	"github.com/root4loot/scope"
)

const (
	AppName = "rescope"
	Version = "2.0.0"
)

func init() {
	log.Init(AppName)
}

type CLI struct {
	Concurrency     int
	Targets         []string
	IncludeList     string
	ExcludeList     string
	TokenBugCrowd   string
	TokenHackerOne  string
	TokenIntigriti  string
	TokenYesWeHack  string
	OutputFile      string
	OutputText      bool
	OutputBurp      bool
	OutputZap       bool
	OutputJson      bool
	OutputJsonLines bool
	ExpandIPRanges  bool
	Debug           bool
}

const usage = `
Usage:
  rescope [options] [<BugBountyURL>...] [-iL <file>] [-eL <file>]

INPUT:
  -iL, --include-list         file containing list of URLs or custom in-scope definitions (newline separated)
  -eL, --exclude-list         file containing list of URLs or custom out-of-scope definitions (newline separated)

OUTPUT:
  -oF, --output-file          output to given file (default: stdout)

OUTPUT FORMAT:
  -oT, --output-text          output simple text (default)
  -oB, --output-burp          output Burp Suite Scope (JSON)
  -oZ, --output-zap           output ZAP Scope (XML)
  -oJ, --output-json          output JSON
  -oJL, --output-json-lines   output JSON lines

OUTPUT FILTER:
  --filter-expand-ip-ranges   output individual IPs instead of IP ranges / CIDRs

AUTHORIZATION:
  --auth-bugcrowd             bugcrowd secret    (_bugcrowd_session=cookie.value) [Optional]
  --auth-hackenproof          hackenproof secret (_hackenproof_session=cookie.value) [Optional]
  --auth-hackerone            hackerone secret   (Authorization bearer token) [Optional]
  --auth-yeswehack            yeswehack secret   (Authorization bearer token) [Optional]
  --auth-intigriti            intigriti secret   (see https://app.intigriti.com/researcher/personal-access-tokens) [Optional]

GENERAL:
  -c, --concurrency           maximum number of concurrent requests (default: 5)
      --debug                 enable debug mode
      --version               display version
`

func parseCLI() ([]string, *CLI, error) {
	var version, help bool
	cli := CLI{}

	flag.StringVar(&cli.IncludeList, "iL", "", "")
	flag.StringVar(&cli.IncludeList, "include-list", "", "")
	flag.StringVar(&cli.ExcludeList, "eL", "", "")
	flag.StringVar(&cli.ExcludeList, "exclude-list", "", "")
	flag.StringVar(&cli.TokenHackerOne, "auth-hackerone", "", "")
	flag.StringVar(&cli.TokenIntigriti, "auth-intigriti", "", "")
	flag.StringVar(&cli.TokenYesWeHack, "auth-yeswehack", "", "")
	flag.StringVar(&cli.TokenBugCrowd, "auth-bugcrowd", "", "")
	flag.StringVar(&cli.TokenBugCrowd, "auth-hackenproof", "", "")
	flag.StringVar(&cli.OutputFile, "oF", "", "")
	flag.StringVar(&cli.OutputFile, "output-file", "", "")
	flag.BoolVar(&cli.OutputText, "oT", false, "")
	flag.BoolVar(&cli.OutputText, "output-text", false, "")
	flag.BoolVar(&cli.OutputBurp, "oB", false, "")
	flag.BoolVar(&cli.OutputBurp, "output-burp", false, "")
	flag.BoolVar(&cli.OutputZap, "oZ", false, "")
	flag.BoolVar(&cli.OutputZap, "output-zap", false, "")
	flag.BoolVar(&cli.OutputJson, "oJ", false, "")
	flag.BoolVar(&cli.OutputJson, "output-json", false, "")
	flag.BoolVar(&cli.OutputJsonLines, "oJL", false, "")
	flag.BoolVar(&cli.OutputJsonLines, "output-json-lines", false, "")
	flag.BoolVar(&cli.ExpandIPRanges, "expand-ip-ranges", false, "")
	flag.IntVar(&cli.Concurrency, "concurrency", 5, "")
	flag.BoolVar(&cli.Debug, "debug", false, "")
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "")
	flag.BoolVar(&version, "version", false, "")

	flag.Parse()

	if version {
		log.Info(Version)
		return nil, nil, fmt.Errorf("version: %s", Version)
	}

	var targets []string
	args := flag.Args()
	if len(args) > 0 {
		targets = args
	}

	if help || version || (len(flag.Args()) == 0 && cli.IncludeList == "" && cli.ExcludeList == "" && !hasStdin()) {
		if help {
			fmt.Fprint(os.Stdout, usage)
			return nil, nil, nil
		} else if version {
			log.Info(Version)
			return nil, nil, fmt.Errorf("version: %s", Version)
		} else {
			return nil, nil, fmt.Errorf("missing URL/file input. No targets provided")
		}
	}

	return targets, &cli, nil
}

func main() {
	args, cli, err := parseCLI()
	if err != nil {
		if err.Error() == fmt.Sprintf("version: %s", Version) {
			log.Info(err.Error())
		} else {
			log.Error(err.Error())
		}
		return
	}
	opts := rescope.DefaultOptions()

	fileIncludes, fileExcludes, err := cli.getInputFileContents()
	if err != nil {
		log.Error("Failed to read input files", "error", err)
		return
	}

	bugBountyURLs := []string{}
	scope := scope.NewScope()

	processScopeList := func(list []string, isInclude bool) {
		for _, item := range list {
			if urlutil.IsURL(item) {
				if rescope.IsBugBountyURL(item) {
					bugBountyURLs = append(bugBountyURLs, item)
				} else {
					if isInclude {
						scope.AddInclude(item)
					} else {
						scope.AddExclude(item)
					}
				}
			} else {
				if isInclude {
					scope.AddInclude(item)
				} else {
					scope.AddExclude(item)
				}
			}
		}
	}

	processScopeList(fileIncludes, true)
	processScopeList(fileExcludes, false)

	if hasStdin() {
		processScopeList(processStdin(), true)
	}
	processScopeList(args, true)

	combinedResult := common.Result{}
	firstResult := true

	if len(fileIncludes) > 0 || len(fileExcludes) > 0 {
		combinedResult = processFileInputs(fileIncludes, fileExcludes, scope, cli)
	}

	cli.setAuthTokens(opts)
	processURLs(bugBountyURLs, opts, cli, scope, &combinedResult, &firstResult)
	printFormattedOutput(&combinedResult, cli)
}

func processFileInputs(fileIncludes, fileExcludes []string, scope *scope.Scope, cli *CLI) common.Result {
	initialResult := &common.Result{
		InScope:  fileIncludes,
		OutScope: fileExcludes,
	}

	scopedResult, err := getScopedResults(*initialResult, *scope)
	if err != nil {
		log.Error("Failed to update results with scope", "error", err)
		return common.Result{}
	}

	scopedResult, err = cli.applyOutputFilters(scopedResult)
	if err != nil {
		log.Error("Failed to apply filters", "error", err)
		return common.Result{}
	}

	return common.Result{
		InScope:  scopedResult.InScope,
		OutScope: scopedResult.OutScope,
	}
}

func processURLs(urls []string, opts *rescope.Options, cli *CLI, scope *scope.Scope, combinedResult *common.Result, firstResult *bool) {
	sem := make(chan struct{}, cli.Concurrency)
	var wg sync.WaitGroup

	for _, url := range urls {
		sem <- struct{}{}
		wg.Add(1)

		go func(url string) {
			defer wg.Done()
			defer func() { <-sem }()

			log.Debug("Processing URL", "url", url)

			_, err := rescope.IdentifyPlatform(url, opts)
			if err != nil {
				log.Error("Unsupported or invalid bug bounty platform", "url", url)
				return
			}

			bugBountyResult, err := rescope.Run(url, opts)
			if err != nil {
				log.Error("Failed to run rescope", "error", err)
				return
			}

			scopedResult, err := getScopedResults(*bugBountyResult, *scope)
			if err != nil {
				log.Error("Failed to update results with scope", "error", err)
				return
			}

			scopedResult, err = cli.applyOutputFilters(scopedResult)
			if err != nil {
				log.Error("Failed to apply filters", "error", err)
				return
			}

			if *firstResult {
				combinedResult.ProgramDetails = scopedResult.ProgramDetails
				*firstResult = false
			}

			combinedResult.InScope = append(combinedResult.InScope, scopedResult.InScope...)
			combinedResult.OutScope = append(combinedResult.OutScope, scopedResult.OutScope...)

		}(url)
	}

	wg.Wait()
}

func hasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	mode := stat.Mode()

	isPipedFromChrDev := (mode & os.ModeCharDevice) == 0
	isPipedFromFIFO := (mode & os.ModeNamedPipe) != 0

	return isPipedFromChrDev || isPipedFromFIFO
}

func processStdin() []string {
	var targets []string
	if hasStdin() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if len(line) > 0 {
				targets = append(targets, strings.Fields(line)...)
			}
		}
	}
	return targets
}

func printFormattedOutput(result *common.Result, cli *CLI) {
	formattedOutput, err := cli.formatOutput(result)
	if err != nil {
		log.Error("Failed to format output", "error", err)
		return
	}

	if cli.OutputFile != "" {
		err := fileutil.WriteStringToFile(cli.OutputFile, formattedOutput)
		if err != nil {
			log.Error("Failed to save output to file", "error", err)
			return
		}
		log.Info("Output saved to file", "file", cli.OutputFile)
	} else {
		fmt.Println(formattedOutput)
	}
}

func getJsonOutput(result *common.Result) (string, error) {
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to serialize Result to JSON: %w", err)
	}
	return string(jsonData), nil
}

func getJsonLineOutput(result *common.Result) (string, error) {
	var lines []string

	addLine := func(inScope, outScope string, programDetails common.BugBountyProgram) error {
		line := map[string]interface{}{
			"program": programDetails,
		}
		if inScope != "" {
			line["in_scope"] = inScope
		}
		if outScope != "" {
			line["out_scope"] = outScope
		}

		data, err := json.Marshal(line)
		if err != nil {
			return fmt.Errorf("failed to serialize to JSON: %w", err)
		}
		lines = append(lines, string(data))
		return nil
	}

	for _, inScope := range result.InScope {
		if err := addLine(inScope, "", result.ProgramDetails); err != nil {
			return "", err
		}
	}

	for _, outScope := range result.OutScope {
		if err := addLine("", outScope, result.ProgramDetails); err != nil {
			return "", err
		}
	}

	return strings.Join(lines, "\n"), nil
}

func getBurpOutput(Result *common.Result) (string, error) {
	var scope config.BurpConfig
	scope.Target.Scope.AdvancedMode = true

	for _, item := range Result.InScope {
		protocol, host, port, file := parseAndReplaceWildcards(item)
		includeEntry := config.BurpInclude{
			Enabled:  true,
			Protocol: protocol,
			Host:     host,
			Port:     port,
			File:     file,
		}

		if protocol == "" {
			includeEntry.Protocol = "any"
		}

		scope.Target.Scope.Include = append(scope.Target.Scope.Include, includeEntry)
	}

	for _, item := range Result.OutScope {
		protocol, host, port, file := parseAndReplaceWildcards(item)
		excludeEntry := config.BurpExclude{
			Enabled:  true,
			Protocol: protocol,
			Host:     host,
			Port:     port,
			File:     file,
		}

		if protocol == "" {
			excludeEntry.Protocol = "any"
		}
		scope.Target.Scope.Exclude = append(scope.Target.Scope.Exclude, excludeEntry)
	}

	output, err := json.MarshalIndent(scope, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to serialize Burp scope to JSON: %w", err)
	}

	return string(output), nil
}

func getZapOutput(Result *common.Result) (string, error) {
	var config config.ZapConfig
	config.Context.Name = "MyContext"
	config.Context.Inscope = "true"

	config.Context.Forceduser = "-1"
	config.Context.Authentication.Type = 0
	config.Context.Authentication.Strategy = "EACH_RESP"
	config.Context.Authentication.Pollurl = ""
	config.Context.Authentication.Polldata = ""
	config.Context.Authentication.Pollfreq = 60
	config.Context.Authentication.Pollunits = "REQUESTS"
	config.Context.Session.Type = 0
	config.Context.Authorization.Type = 0
	config.Context.Authorization.Basic.Header = ""
	config.Context.Authorization.Basic.Body = ""
	config.Context.Authorization.Basic.Logic = "AND"
	config.Context.Authorization.Basic.Code = -1

	appendWildcardIfNeeded := func(protocol, host, port, file string) string {
		if file == "" || !strings.HasSuffix(file, ".*") {
			file += ".*"
		}

		if protocol == "" {
			protocol = "http(s)?"
		}

		if port != "" {
			host = fmt.Sprintf("%s:%s", host, port)
		}

		return fmt.Sprintf("^%s://%s%s$", protocol, host, file)
	}

	processScope := func(scopeItems []string, appendTo *[]string) {
		for _, item := range scopeItems {
			if iputil.IsCIDR(item) || iputil.IsIPRange(item) || iputil.IsIP(item) {
				ips, err := ipRangeToIPs([]string{item}) // Convert IP ranges/CIDRs to individual IPs
				if err != nil {
					continue
				}
				for _, ip := range ips {
					protocol := "http(s)?"
					host, file := ip, ""
					regex := appendWildcardIfNeeded(protocol, host, "", file)
					*appendTo = append(*appendTo, regex)
				}
			} else {
				// For normal URLs
				protocol, host, port, file := parseAndReplaceWildcards(item)
				regex := appendWildcardIfNeeded(protocol, host, port, file)
				*appendTo = append(*appendTo, regex)
			}
		}
	}

	processScope(Result.InScope, &config.Context.Incregexes)
	processScope(Result.OutScope, &config.Context.Excregexes)

	config.Context.Tech.Include = []string{
		"Db", "Db.Firebird", "Db.HypersonicSQL", "Db.IBM DB2", "Db.Microsoft Access", "Db.Microsoft SQL Server",
		"Db.MySQL", "Db.Oracle", "Db.PostgreSQL", "Db.SAP MaxDB", "Db.SQLite", "Db.Sybase",
		"Language", "Language.ASP", "Language.C", "Language.PHP", "Language.XML",
		"OS", "OS.Linux", "OS.MacOS", "OS.Windows",
		"SCM", "SCM.Git", "SCM.SVN",
		"WS", "WS.Apache", "WS.IIS", "WS.Tomcat",
	}

	config.Context.Urlparser.Class = "org.zaproxy.zap.model.StandardParameterParser"
	config.Context.Urlparser.Config = "{\"kvps\":\"&\",\"kvs\":\"=\",\"struct\":[]}"
	config.Context.Postparser.Class = "org.zaproxy.zap.model.StandardParameterParser"
	config.Context.Postparser.Config = "{\"kvps\":\"&\",\"kvs\":\"=\",\"struct\":[]}"

	output, err := xml.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to serialize ZAP scope to XML: %w", err)
	}

	return xml.Header + string(output), nil
}

func getSimpleTextOutput(result *common.Result) string {
	var builder strings.Builder

	for _, scopeItem := range result.InScope {
		builder.WriteString("In-Scope: ")
		builder.WriteString(scopeItem)
		builder.WriteString("\n")
	}

	for _, scopeItem := range result.OutScope {
		builder.WriteString("Out-Scope: ")
		builder.WriteString(scopeItem)
		builder.WriteString("\n")
	}

	return strings.TrimRight(builder.String(), "\n")
}

func ipRangeToIPs(scopeItems []string) ([]string, error) {
	var converted []string
	for _, item := range scopeItems {
		if iputil.IsIPRange(item) {
			ips, err := iputil.ParseIPRange(item)
			if err != nil {
				return nil, err
			}
			for _, ip := range ips {
				converted = append(converted, ip.String())
			}
			continue
		} else if iputil.IsCIDR(item) {
			ips, err := iputil.ParseCIDR(item)
			if err != nil {
				return nil, err
			}
			for _, ip := range ips {
				converted = append(converted, ip.String())
			}
			continue
		}
		converted = append(converted, item)
	}

	return converted, nil
}

func parseAndReplaceWildcards(input string) (protocol, host, port, file string) {
	if iputil.IsIP(input) || iputil.IsCIDR(input) || iputil.IsIPRange(input) {
		return "any", input, "", ""
	}

	r1 := regexp.MustCompile(`(?i)^(?:(?P<protocol>[a-z3]+):\/\/)?(?P<host>(?:(?:[a-z]+\.)?(?:\*\.)?\*?[a-zA-Z0-9-.*]+(?:\.[a-z]+|\.\*)))?(?::(?P<port>\d+))?(?P<file>[\/A-Za-z0-9-._~:/?#@!$&'*+=]*)?$`)

	matches := r1.FindStringSubmatch(input)
	if len(matches) == 0 {
		return "", "", "", ""
	}

	protocol = matches[r1.SubexpIndex("protocol")]
	host = matches[r1.SubexpIndex("host")]
	host = strings.ReplaceAll(host, "*", ".*")
	port = matches[r1.SubexpIndex("port")]
	file = matches[r1.SubexpIndex("file")]
	file = strings.ReplaceAll(file, "*", ".*")

	return
}

func getScopedResults(result common.Result, scope scope.Scope) (*common.Result, error) {
	var newResult common.Result
	newResult.ProgramDetails = result.ProgramDetails

	newResult.InScope = append(newResult.InScope, result.InScope...)
	newResult.OutScope = append(newResult.OutScope, result.OutScope...)

	for _, include := range scope.GetIncludes() {
		if !sliceutil.Contains(newResult.InScope, include) {
			newResult.InScope = append(newResult.InScope, include)
		}
	}

	for _, exclude := range scope.GetExcludes() {
		if !sliceutil.Contains(newResult.OutScope, exclude) {
			newResult.OutScope = append(newResult.OutScope, exclude)
		}
	}

	return &newResult, nil
}

func (cli *CLI) getInputFileContents() (includeTargets, excludeTargets []string, err error) {
	if cli.IncludeList != "" {
		includeTargets, err = fileutil.ReadFile(cli.IncludeList)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read bug bounty URLs file: %w", err)
		}
	}

	if cli.ExcludeList != "" {
		excludeTargets, err = fileutil.ReadFile(cli.ExcludeList)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read include file: %w", err)
		}
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to read exclude file: %w", err)
	}

	return includeTargets, excludeTargets, nil
}

func (cli *CLI) formatOutput(result *common.Result) (string, error) {
	switch {
	case cli.OutputJson:
		return getJsonOutput(result)
	case cli.OutputJsonLines:
		return getJsonLineOutput(result)
	case cli.OutputBurp:
		return getBurpOutput(result)
	case cli.OutputZap:
		return getZapOutput(result)
	default:
		return getSimpleTextOutput(result), nil
	}
}

func (cli *CLI) applyOutputFilters(Result *common.Result) (*common.Result, error) {
	if cli.ExpandIPRanges {
		var err error
		Result.InScope, err = ipRangeToIPs(Result.InScope)
		if err != nil {
			return Result, fmt.Errorf("failed to convert IP ranges: %w", err)
		}

		Result.OutScope, err = ipRangeToIPs(Result.OutScope)
		if err != nil {
			return Result, fmt.Errorf("failed to convert IP ranges: %w", err)
		}
	}

	return Result, nil
}

func (cli *CLI) setAuthTokens(opts *rescope.Options) {
	opts.AuthHackerOne = cli.TokenHackerOne
	opts.AuthIntigriti = cli.TokenIntigriti
	opts.AuthYesWeHack = cli.TokenYesWeHack
	opts.AuthBugcrowd = cli.TokenBugCrowd
	opts.AuthHackenProof = cli.TokenBugCrowd

	if cli.Debug {
		opts.Debug = true
	}
}
