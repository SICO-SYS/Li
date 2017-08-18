/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"github.com/getsentry/raven-go"
	"golang.org/x/net/context"

	"github.com/SiCo-Ops/Li/controller/qcloud"
	"github.com/SiCo-Ops/Pb"
)

var (
	requestUrl, requestParamString, signature string
)

func (q *CloudAPIService) RequestRPC(ctx context.Context, in *pb.CloudAPICall) (*pb.CloudAPIBack, error) {
	switch in.Cloud {
	case "qcloud":
		requestUrl = in.Service + ".api.qcloud.com/v2/index.php"
		requestParamString = qcloud.SignatureString(in.Service, in.Action, in.Region, in.CloudId, in.Params)
		signature = qcloud.Signature(requestUrl, requestParamString, in.CloudKey)
		res, err := qcloud.Request(requestUrl, requestParamString, signature)
		if err != nil {
			raven.CaptureError(err, nil)
			return &pb.CloudAPIBack{Code: 1, Msg: "qcloud maybe probl"}, nil
		}
		return &pb.CloudAPIBack{Code: 0, Msg: "Success", Data: res}, nil
	default:
		return &pb.CloudAPIBack{Code: 2, Msg: "Cloud Not support yet"}, nil
	}
}
