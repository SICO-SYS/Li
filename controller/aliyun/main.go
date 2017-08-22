/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package aliyun

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/SiCo-Ops/public"
)

func URLEncode(signstr string) string {
	result := url.QueryEscape(signstr)
	result = strings.Replace(result, "+", "%20", -1)
	result = strings.Replace(result, "*", "%2A", -1)
	result = strings.Replace(result, "%7E", "~", -1)
	return result
}

func SignatureString(service, action, region, secretId string, extraParams map[string]string) string {
	params := make(map[string]string)
	var sortparams = []string{}
	params["Format"] = "JSON"
	sortparams = append(sortparams, "Format")
	params["AccessKeyId"] = secretId
	sortparams = append(sortparams, "AccessKeyId")
	params["SignatureMethod"] = "HMAC-SHA1"
	sortparams = append(sortparams, "SignatureMethod")
	params["Timestamp"] = public.CurrentUTCISO8601()
	sortparams = append(sortparams, "Timestamp")
	params["SignatureVersion"] = "1.0"
	sortparams = append(sortparams, "SignatureVersion")
	params["SignatureNonce"] = public.GenerateNonce()
	sortparams = append(sortparams, "SignatureNonce")
	params["Action"] = action
	sortparams = append(sortparams, "Action")
	params["RegionId"] = region
	sortparams = append(sortparams, "RegionId")

	switch service {
	case "ecs":
		params["Version"] = "2014-05-26"
		sortparams = append(sortparams, "Version")
	}

	for paramKey, paramValue := range extraParams {
		params[paramKey] = paramValue
		sortparams = append(sortparams, paramKey)
	}
	sort.Strings(sortparams)
	requestParamString := ""
	var paramstr = []string{}
	for _, request_key := range sortparams {
		paramstr = append(paramstr, URLEncode(request_key)+"="+URLEncode(params[request_key]))
	}
	requestParamString += strings.Join(paramstr, "&")

	return requestParamString

}

func Signature(requestParamString, secretKey string) string {
	signstr := "GET&%2F&" + URLEncode(requestParamString)
	signatrue := public.Hmac1ToBase64(secretKey+"&", signstr, true)
	return signatrue
}

func URL(scheme, service, region string) string {
	switch service {
	default:
		return scheme + service + ".aliyuncs.com"
	}
}

func Request(requestUrl, requestParamString, signature string) ([]byte, error) {
	resp, err := http.Get(requestUrl + "/?" + requestParamString + "&Signature=" + signature)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(resp.Body)
	return res, err
}
