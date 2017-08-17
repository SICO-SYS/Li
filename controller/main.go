/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"github.com/getsentry/raven-go"
	"google.golang.org/grpc"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/cfg"
	"github.com/SiCo-Ops/dao/mongo"
)

var (
	config    = cfg.Config
	RPCServer = grpc.NewServer()
)

type CloudAPIService struct{}
type CloudTokenService struct{}

func init() {
	defer func() {
		recover()
	}()
	mongo.CloudEnsureIndexes()
	pb.RegisterCloudAPIServiceServer(RPCServer, &CloudAPIService{})
	pb.RegisterCloudTokenServiceServer(RPCServer, &CloudTokenService{})

	if config.Sentry.Enable {
		raven.SetDSN(config.Sentry.DSN)
	}
}
