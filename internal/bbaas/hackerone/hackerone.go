//
// Author: Daniel Antonsen (@root4loot)
// Distributed Under MIT License
//

package hackerone

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"unsafe"

	doerror "github.com/root4loot/rescope/internal/bbaas/pkg/errors"
	"github.com/root4loot/rescope/internal/bbaas/pkg/request"
)

var scope []string
var isAssetURL bool

// Scrape tries to grab scope table for a given program on hackerone.com
func Scrape(url string) string {
	re := regexp.MustCompile(`([\w-]+)\/([\w-]+$)`)
	match := re.FindStringSubmatch(url)
	program := strings.ToLower(match[2])
	endpoint := "https://hackerone.com/graphql?"

	// clear global slice
	scope = nil

	var data = []byte(`{  
		"query":"query Team_assets($first_0:Int!) {query {id,...F0}} fragment F0 on Query {_teamAgUhl:team(handle:\"` + program + `\") {handle,_structured_scope_versions2ZWKHQ:structured_scope_versions(archived:false) {max_updated_at},_structured_scopeszxYtW:structured_scopes(first:$first_0,archived:false,eligible_for_submission:true) {edges {node {asset_type, asset_identifier}},pageInfo {hasNextPage,hasPreviousPage}},_structured_scopes3FF98f:structured_scopes(first:$first_0,archived:false,eligible_for_submission:false) {edges {node {asset_type,asset_identifier,},},},},}",
		"variables":{  
		   "first_0":1337
		}
	 }`)

	// GET to check if program is reachable
	_, responseCode := (request.GET(url))
	if responseCode != 200 {
		doerror.NoScope(url)
	}

	resB, _ := (Post(endpoint, data))
	resS := BytesToString(resB)

	re = regexp.MustCompile(`\"edges":\[(.*?)\]`)
	scopeSplit := re.FindAllString(resS, -1)
	re = regexp.MustCompile(`asset_type":"(URL|CIDR|IP|IP-RANGE|RANGE)","asset_identifier":"(.*?)"`)

	inScope := re.FindAllStringSubmatch(scopeSplit[0], -1)
	outScope := re.FindAllStringSubmatch(scopeSplit[1], -1)

	scope = append(scope, "!INCLUDE")
	for _, m := range inScope {
		scope = append(scope, m[2])
	}

	scope = append(scope, "!EXCLUDE")
	for _, m := range outScope {
		scope = append(scope, m[2])
	}

	return strings.Join(scope, "\n")
}

// BytesToString converts byte array to string
func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

// Post makes a post request with custom X-Auth-Token header
func Post(url string, data []byte) ([]byte, int) {
	token := os.Getenv("H1_TOKEN")
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")
	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("Cookie", "__Host-session=SGFVODkyZDFJaU5pUXB3WFFpTTJzN2J4SU52MkU4Nmg3Nmo3dVo3YkJlQ29INEtVSGFoZW44cWlSWFdLLzBGR28xQTBNUTArMkVkaSt0bFVYV3loS2RvbTJqZnprK1JqWVJTS0gxNGxuMjBxUVJzUzRkeHdNT2VtSmc1dzQ4aXdQS0lid0FsaStrcS9pRDQ2WE9sclNEMFBvTnVPMGorQUkwVDE4NXVxcUhLOHd5VDFsWndQZUVYOXh3Q05CWURaU3NkODRVUFJKZGNTZE1wa3YxRlB5ZW5FbTF2cTlkeVdoTWw2NHprejFjVnVnY24wMWlQdkdOMWJNZGszdWFRckJaazI5QWh6amNPMGFaWVVZWlU5ODB1UmVwa3RDTVFzbG8yWHRvWGpqTzA9LS1SVkxyRGk2YXlDSmtIMmV1MXY2RitnPT0%3D--d4d79937e04894b7c2cabbf0ce56ed69ae3cdd8f;")
	req.Header.Set("X-Csrf-Token", "nonlX8VXBB5caid8i5u5njLdgLJKYjn/0sFd0A9+HWZz/6M4Bjb9iu16Ac/7n0mkvk3jAywTyTKvx+S2GnUbJQ==")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	respS := resp.StatusCode

	// check response
	if err != nil {
		doerror.NoResponse(url)
	}

	// close response
	defer resp.Body.Close()

	// JSON response body
	respB, _ := ioutil.ReadAll(resp.Body)

	return respB, respS
}
