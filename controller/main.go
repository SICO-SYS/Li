/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	"github.com/getsentry/raven-go"
	"google.golang.org/grpc"
	"io/ioutil"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/dao/mongo"
)

var (
	RPCServer = grpc.NewServer()
	getAction map[string]interface{}
)

type CloudAPIService struct{}

func mapAction(cloud string, service string, action string) (string, bool) {
	d, err := ioutil.ReadFile("action.json")
	if err != nil {
		raven.CaptureError(err, nil)
	}
	json.Unmarshal(d, &getAction)

	getCloud, ok := getAction[action].(map[string]interface{})
	if ok {
		getService, ok := getCloud[cloud].(map[string]interface{})
		if ok {
			value, ok := getService[service].(string)
			if ok {
				return value, true
			}
			return "Service unsupported", false
		}
		return "Cloud unsupported", false
	}
	return "Action unsupported", false
}

func init() {
	defer func() {
		recover()
	}()
	mongo.CloudEnsureIndexes()
	pb.RegisterCloudAPIServiceServer(RPCServer, &CloudAPIService{})

	if config.Sentry.Enable {
		raven.SetDSN(config.Sentry.DSN)
	}
}
