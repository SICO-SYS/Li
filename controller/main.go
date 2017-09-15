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
	"github.com/SiCo-Ops/cfg/v2"
	"github.com/SiCo-Ops/dao/mongo"
	"github.com/SiCo-Ops/dao/redis"
)

var (
	config              cfg.ConfigItems
	configPool          = redis.Pool("", "", "")
	RPCServer           = grpc.NewServer()
	cloudDB, cloudDBErr = mongo.Dial("", "", "")
)

type CloudAPIService struct{}
type CloudTokenService struct{}

func ServePort() string {
	return config.RpcLiPort
}

func init() {
	defer func() {
		recover()
	}()
	data := cfg.ReadLocalFile()

	if data != nil {
		cfg.Unmarshal(data, &config)
	}

	configPool = redis.Pool(config.RedisConfigHost, config.RedisConfigPort, config.RedisConfigAuth)
	configs, _ := redis.Hgetall(configPool, "system.config")
	cfg.Map2struct(configs, &config)

	cloudDB, cloudDBErr = mongo.Dial(config.MongoCloudAddress, config.MongoCloudUsername, config.MongoCloudPassword)

	mongo.CloudEnsureIndexes(cloudDB)
	pb.RegisterCloudAPIServiceServer(RPCServer, &CloudAPIService{})
	pb.RegisterCloudTokenServiceServer(RPCServer, &CloudTokenService{})

	if config.SentryLiStatus == "active" && config.SentryLiDSN != "" {
		raven.SetDSN(config.SentryLiDSN)
	}
}
