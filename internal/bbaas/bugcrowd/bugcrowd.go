//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2021 root4loot
//

package bugcrowd

import (
	"fmt"
	"regexp"
	"strings"

	errors "github.com/root4loot/rescope/internal/bbaas/pkg/errors"
	req "github.com/root4loot/rescope/internal/bbaas/pkg/request"
)

// Scrape returns a string containing scope that was scraped from the given program on bugcrowd.com
func Scrape(url string) string {
	var scope []string
	re := regexp.MustCompile(`([\w-]+)\/([\w-]+$)`)
	re_nameKey := regexp.MustCompile(`name":"(.*?)"`)
	match := re.FindStringSubmatch(url)
	scheme := "https://"
	domain := "bugcrowd.com/"
	program := match[2]
	target_groups := scheme + domain + program + "/target_groups"

	// GET target groups
	resp, status := req.GET(target_groups)
	if status != 200 {
		errors.BadStatusCode(url, status)
	}

	re_inScopeTargetGroup := regexp.MustCompile(`in_scope":true(.*?)"targets_url":"(.*?)"`)
	re_OutScopeTargetGroup := regexp.MustCompile(`in_scope":false(.*?)"targets_url":"(.*?)"`)
	match_InScopeTargetGroup := re_inScopeTargetGroup.FindAllStringSubmatch(resp, -1)
	match_OutScopeTargetgroup := re_OutScopeTargetGroup.FindAllStringSubmatch(resp, -1)

	if match_InScopeTargetGroup != nil {
		scope = append(scope, "!INCLUDE")
		for _, match := range match_InScopeTargetGroup {
			resp, _ := req.GET(scheme + domain + match[2])
			names := re_nameKey.FindAllStringSubmatch(resp, -1)
			for _, v := range names {
				fmt.Println(v[1])
				scope = append(scope, v[1])
			}
		}
	}

	if match_OutScopeTargetgroup != nil {
		scope = append(scope, "!EXCLUDE")
		for _, match := range match_OutScopeTargetgroup {
			resp, _ := req.GET(scheme + domain + match[2])
			names := re_nameKey.FindAllStringSubmatch(resp, -1)
			for _, v := range names {
				fmt.Println(v[1])
				scope = append(scope, v[1])
			}
		}
	}

	return strings.Join(scope, "\n")
}
