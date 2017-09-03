/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
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
	action := in.Action
	service := in.Service
	region := in.Region
	secretId := in.CloudId
	secretKey := in.CloudKey
	extraParams := in.Params
	switch in.Cloud {
	case "qcloud":
		requestUrl = qcloudSDK.Host(service, region)
		requestParamString = qcloudSDK.SignatureString(service, action, region, secretId, extraParams)
		signature = qcloudSDK.Signature(requestUrl, requestParamString, secretKey)
		res, err := qcloudSDK.Request(requestUrl, requestParamString, signature)
		if err != nil {
			raven.CaptureError(err, nil)
			return &pb.CloudAPIBack{Code: 1, Msg: "qcloud maybe probl"}, nil
		}
		return &pb.CloudAPIBack{Code: 0, Msg: "Success", Data: res}, nil
	case "aliyun":
		requestUrl = aliyunSDK.URL("http://", service, region)
		requestParamString = aliyunSDK.SignatureString(service, action, region, secretId, extraParams)
		signature = aliyunSDK.Signature(requestParamString, secretKey)
		res, err := aliyunSDK.Request(requestUrl, requestParamString, signature)
		if err != nil {
			raven.CaptureError(err, nil)
			return &pb.CloudAPIBack{Code: 2, Msg: "Aliyun maybe probl"}, nil
		}
		return &pb.CloudAPIBack{Code: 0, Msg: "Success", Data: res}, nil
	case "aws":
		if region == "" {
			region = "us-east-1"
		}
		amzDate, amzDatetime := awsSDK.Dates()

		host := awsSDK.Host(service, region)
		requestUrl = "https://" + host
		credentialScope := awsSDK.CredentialScope(amzDate, service, region)
		requestParamString = awsSDK.CanonicalQueryString(service, action, credentialScope, secretId, amzDatetime, extraParams)
		signature = awsSDK.Signature(awsSDK.SignatureString(credentialScope, awsSDK.CanonicalRequest(requestParamString, awsSDK.CanonicalHost(host)), amzDatetime), awsSDK.SignatureKey(secretKey, region, service, amzDate))
		res, err := awsSDK.Request(requestUrl, requestParamString, signature)
		if err != nil {
			raven.CaptureError(err, nil)
			return &pb.CloudAPIBack{Code: 2, Msg: "AWS maybe probl"}, nil
		}
		return &pb.CloudAPIBack{Code: 0, Msg: "Success", Data: res}, nil
	default:
		return &pb.CloudAPIBack{Code: 2, Msg: "Cloud Not support yet"}, nil
	}
}
