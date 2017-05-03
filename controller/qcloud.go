/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/SiCo-DevOps/Pb"
	"github.com/SiCo-DevOps/public"
)

func (q *Cloud_API) Qcloud(ctx context.Context, in *pb.CloudRequest) (*pb.CloudResponse, error) {
	cloud_action, ok := transAction("qcloud", in.Bsns, in.Action)
	if !ok {
		return &pb.CloudResponse{Code: 2, Msg: cloud_action}, nil
	}
	// scheme := "https://"
	host := in.Bsns + ".api.qcloud.com/v2/index.php"
	params := make(map[string]string)
	var sortparams = []string{}
	params["Action"] = cloud_action
	sortparams = append(sortparams, "Action")
	params["Nonce"] = public.GenNonce()
	sortparams = append(sortparams, "Nonce")
	params["Region"] = in.Region
	sortparams = append(sortparams, "Region")
	params["Timestamp"] = public.TS()
	sortparams = append(sortparams, "Timestamp")
	params["SecretId"] = in.CloudId
	sortparams = append(sortparams, "SecretId")
	params["SignatureMethod"] = "HmacSHA256"
	sortparams = append(sortparams, "SignatureMethod")
	for _, param_value := range in.Params {
		params[param_value.Key] = param_value.Value
		sortparams = append(sortparams, param_value.Key)
	}
	sort.Strings(sortparams)
	requeststr_ask := "?"
	requeststr := ""
	var paramstr = []string{}
	for _, request_key := range sortparams {
		paramstr = append(paramstr, request_key+"="+params[request_key])
	}
	requeststr += strings.Join(paramstr, "&")
	signstr := "POST" + host + requeststr_ask + requeststr
	signatrue := public.Hmac256ToBase64(in.CloudKey, signstr, true)

	resp, _ := http.Post("https://"+host, "application/x-www-form-urlencoded", strings.NewReader(requeststr+"&Signature="+signatrue))
	defer resp.Body.Close()
	res, _ := ioutil.ReadAll(resp.Body)
	var v json.RawMessage
	json.Unmarshal(res, &v)
	return &pb.CloudResponse{Code: 0, Msg: "Success", Data: v}, nil
}
