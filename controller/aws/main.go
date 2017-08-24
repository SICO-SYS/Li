/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package aws

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/SiCo-Ops/public"
)

var (
	amzDate      string = public.CurrentYYYMMDD()
	amzDatetime  string = public.CurrentYYYMMDDTHHMMSSZ()
	amzAlgorithm string = "AWS4-HMAC-SHA256"
)

func URLEncode(signstr string) string {
	result := url.QueryEscape(signstr)
	result = strings.Replace(result, "+", "%20", -1)
	return result
}

func Host(service, region string) string {
	switch service {
	case "devpay":
		return "ls.amazonaws.com"
	case "iam":
		return "iam.amazonaws.com"
	case "s3":
		switch region {
		case "us-east-1":
			return "s3.amazonaws.com"
		default:
			return service + "-" + region + ".amazonaws.com"
		}
	default:
		return service + "." + region + ".amazonaws.com"
	}
}

func CredentialScope(service, region string) string {
	return strings.Join([]string{amzDate, region, service, "aws4_request"}, "/")
}

func CanonicalQueryString(service, action, region, secretId string, extraParams map[string]string) string {
	params := make(map[string]string)
	var sortparams = []string{}
	params["Action"] = action
	sortparams = append(sortparams, "Action")
	params["X-Amz-Algorithm"] = amzAlgorithm
	sortparams = append(sortparams, "X-Amz-Algorithm")
	params["X-Amz-Credential"] = secretId + "/" + CredentialScope(service, region)
	sortparams = append(sortparams, "X-Amz-Credential")
	params["X-Amz-Date"] = amzDatetime
	sortparams = append(sortparams, "X-Amz-Date")
	params["X-Amz-SignedHeaders"] = "host"
	sortparams = append(sortparams, "X-Amz-SignedHeaders")

	switch service {
	case "s3":
		params["Version"] = "2006-03-01"
		sortparams = append(sortparams, "Version")
	case "elasticache":
		params["Version"] = "2015-02-02"
		sortparams = append(sortparams, "Version")
	default:
		params["Version"] = "2016-11-15"
		sortparams = append(sortparams, "Version")
	}

	for paramKey, paramValue := range extraParams {
		params[paramKey] = paramValue
		sortparams = append(sortparams, paramKey)
	}
	sort.Strings(sortparams)
	var paramstr = []string{}
	for _, request_key := range sortparams {
		paramstr = append(paramstr, URLEncode(request_key)+"="+URLEncode(params[request_key]))
	}
	return strings.Join(paramstr, "&")
}

func CanonicalHost(host string) string {
	return "host:" + host + "\n"
}

func CanonicalRequest(queryString, canonicalhost string) string {
	return strings.Join([]string{"GET", "/", queryString, canonicalhost, "host", public.EncryptWithSha256("")}, "\n")
}

func SignatureString(credentialScope, canonicalRequest string) string {
	return strings.Join([]string{amzAlgorithm, amzDatetime, credentialScope, public.EncryptWithSha256(canonicalRequest)}, "\n")
}

func SignatureKey(secretKey, region, service string) []byte {
	key := public.EncryptBytesWithHmacSha256([]byte("AWS4"+secretKey), amzDate)
	key = public.EncryptBytesWithHmacSha256(key, region)
	key = public.EncryptBytesWithHmacSha256(key, service)
	key = public.EncryptBytesWithHmacSha256(key, "aws4_request")
	return key
}

func Signature(signatureString string, signatureKey []byte) string {
	signatrue := public.EncryptBytesWithHmacSha256(signatureKey, signatureString)
	return public.EncodingBytesToHex(signatrue)
}

func Request(requestUrl, requestParamString, signature string) ([]byte, error) {
	resp, err := http.Get(requestUrl + "/?" + requestParamString + "&X-Amz-Signature=" + signature)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(resp.Body)
	return res, err
}
