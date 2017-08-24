/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package qcloud

import (
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/SiCo-Ops/public"
)

func SignatureString(service, action, region, secretId string, extraParams map[string]string) string {
	params := make(map[string]string)
	var sortparams = []string{}
	params["Action"] = action
	sortparams = append(sortparams, "Action")
	params["Nonce"] = public.GenerateNonce()
	sortparams = append(sortparams, "Nonce")
	params["Region"] = region
	sortparams = append(sortparams, "Region")
	params["Timestamp"] = public.CurrentTimeStamp()
	sortparams = append(sortparams, "Timestamp")
	params["SecretId"] = secretId
	sortparams = append(sortparams, "SecretId")
	params["SignatureMethod"] = "HmacSHA256"
	sortparams = append(sortparams, "SignatureMethod")

	switch service {
	case "cvm":
		params["Version"] = "2017-03-12"
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
		paramstr = append(paramstr, request_key+"="+params[request_key])
	}
	requestParamString += strings.Join(paramstr, "&")

	return requestParamString

}

func Signature(requestUrl, requestParamString, secretKey string) string {
	signstr := "POST" + requestUrl + "?" + requestParamString
	signatrue := public.Hmac256ToBase64(secretKey, signstr, true)
	return signatrue
}

func Host(service, region string) string {
	switch service {
	default:
		return service + ".api.qcloud.com/v2/index.php"
	}
}

func Request(requestUrl, requestParamString, signature string) ([]byte, error) {
	resp, err := http.Post("https://"+requestUrl, "application/x-www-form-urlencoded", strings.NewReader(requestParamString+"&Signature="+signature))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(resp.Body)
	return res, err
}
