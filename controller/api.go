/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/xml"
	"fmt"
	"github.com/getsentry/raven-go"
	"golang.org/x/net/context"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/cloud-go-sdk/aliyun"
	"github.com/SiCo-Ops/cloud-go-sdk/aws"
	"github.com/SiCo-Ops/cloud-go-sdk/qcloud"
)

var (
	requestUrl, requestParamString, signature string
)

func (q *CloudAPIService) RequestRPC(ctx context.Context, in *pb.CloudAPICall) (*pb.CloudAPIBack, error) {
	defer func() {
		recover()
	}()
	switch in.Cloud {
	case "qcloud":
		requestUrl = qcloudSDK.Host(in.Service, in.Region)
		requestParamString = qcloudSDK.SignatureString(in.Service, in.Action, in.Region, in.CloudId, in.Params)
		signature = qcloudSDK.Signature(requestUrl, requestParamString, in.CloudKey)
		res, err := qcloudSDK.Request(requestUrl, requestParamString, signature)
		if err != nil {
			raven.CaptureError(err, nil)
			return &pb.CloudAPIBack{Code: 1, Msg: "qcloud maybe probl"}, nil
		}
		return &pb.CloudAPIBack{Code: 0, Msg: "Success", Data: res}, nil
	case "aliyun":
		requestUrl = aliyunSDK.URL("http://", in.Service, in.Region)
		requestParamString = aliyunSDK.SignatureString(in.Service, in.Action, in.Region, in.CloudId, in.Params)
		signature = aliyunSDK.Signature(requestParamString, in.CloudKey)
		res, err := aliyunSDK.Request(requestUrl, requestParamString, signature)
		if err != nil {
			raven.CaptureError(err, nil)
			return &pb.CloudAPIBack{Code: 2, Msg: "Aliyun maybe probl"}, nil
		}
		return &pb.CloudAPIBack{Code: 0, Msg: "Success", Data: res}, nil
	case "aws":
		action := in.Action
		service := in.Service
		region := in.Region
		secretId := in.CloudId
		secretKey := in.CloudKey
		extraParams := in.Params
		if region == "" {
			region = "us-east-1"
		}
		host := awsSDK.Host(service, region)
		requestUrl = "https://" + host
		requestParamString = awsSDK.CanonicalQueryString(service, action, region, secretId, extraParams)
		signature = awsSDK.Signature(awsSDK.SignatureString(awsSDK.CredentialScope(service, region), awsSDK.CanonicalRequest(requestParamString, awsSDK.CanonicalHost(host))), awsSDK.SignatureKey(secretKey, region, service))
		res, err := awsSDK.Request(requestUrl, requestParamString, signature)
		var v interface{}
		xml.Unmarshal(res, &v)
		fmt.Println(v)
		if err != nil {
			raven.CaptureError(err, nil)
			return &pb.CloudAPIBack{Code: 2, Msg: "AWS maybe probl"}, nil
		}
		return &pb.CloudAPIBack{Code: 0, Msg: "Success", Data: res}, nil
	default:
		return &pb.CloudAPIBack{Code: 2, Msg: "Cloud Not support yet"}, nil
	}
}
