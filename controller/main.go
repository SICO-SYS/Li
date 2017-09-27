/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"github.com/getsentry/raven-go"
	"google.golang.org/grpc"
	"log"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/cfg/v2"
	"github.com/SiCo-Ops/dao/mongo"
	"github.com/SiCo-Ops/dao/redis"
)

var (
	config              cfg.ConfigItems
	configPool          = redis.NewPool()
	RPCServer           = grpc.NewServer()
	cloudDB, cloudDBErr = mongo.NewDial()
)

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

	configPool = redis.InitPool(config.RedisConfigHost, config.RedisConfigPort, config.RedisConfigAuth)
	configs, err := redis.Hgetall(configPool, "system.config")
	if err != nil {
		log.Fatalln(err)
	}
	cfg.Map2struct(configs, &config)

	cloudDB, cloudDBErr = mongo.InitDial(config.MongoCloudAddress, config.MongoCloudUsername, config.MongoCloudPassword)
	if cloudDBErr != nil {
		log.Fatalln(cloudDBErr)
	}
	err = mongo.CloudEnsureIndexes(cloudDB)
	if err != nil {
		log.Fatalln(err)
	}

	pb.RegisterCloudAPIServiceServer(RPCServer, &CloudAPIService{})
	pb.RegisterCloudTokenServiceServer(RPCServer, &CloudTokenService{})

	if config.SentryLiStatus == "active" && config.SentryLiDSN != "" {
		raven.SetDSN(config.SentryLiDSN)
	}
}
