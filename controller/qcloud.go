/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/public"
)

func (q *CloudAPIService) QcloudRPC(ctx context.Context, in *pb.CloudAPICall) (*pb.CloudAPIBack, error) {
	action, ok := mapAction("qcloud", in.Service, in.Action)
	if !ok {
		return &pb.CloudAPIBack{Code: 2, Msg: action}, nil
	}
	scheme := "https://"
	host := in.Service + ".api.qcloud.com"
	path := "/v2/index.php"
	params := make(map[string]string)
	var sortparams = []string{}
	params["Action"] = action
	sortparams = append(sortparams, "Action")
	params["Nonce"] = public.GenerateNonce()
	sortparams = append(sortparams, "Nonce")
	params["Region"] = in.Region
	sortparams = append(sortparams, "Region")
	params["Timestamp"] = public.CurrentTimeStamp()
	sortparams = append(sortparams, "Timestamp")
	params["SecretId"] = in.CloudId
	sortparams = append(sortparams, "SecretId")
	params["SignatureMethod"] = "HmacSHA256"
	sortparams = append(sortparams, "SignatureMethod")

	if in.Service == "cvm" {
		params["Version"] = "2017-03-12"
		sortparams = append(sortparams, "Version")

	}

	for paramKey, paramValue := range in.Params {
		params[paramKey] = paramValue
		sortparams = append(sortparams, paramKey)
	}
	sort.Strings(sortparams)
	requeststr := ""
	var paramstr = []string{}
	for _, request_key := range sortparams {
		paramstr = append(paramstr, request_key+"="+params[request_key])
	}
	requeststr += strings.Join(paramstr, "&")
	signstr := "POST" + host + path + "?" + requeststr
	signatrue := public.Hmac256ToBase64(in.CloudKey, signstr, true)

	resp, _ := http.Post(scheme+host+path, "application/x-www-form-urlencoded", strings.NewReader(requeststr+"&Signature="+signatrue))
	defer resp.Body.Close()
	res, _ := ioutil.ReadAll(resp.Body)
	return &pb.CloudAPIBack{Code: 0, Msg: "Success", Data: res}, nil
}
