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
	"github.com/SiCo-Ops/cfg"
	"github.com/SiCo-Ops/dao/mongo"
)

const (
	configPath string = "config.json"
)

var (
	config              cfg.ConfigItems
	RPCServer           = grpc.NewServer()
	cloudDB, cloudDBErr = mongo.NewDial()
)

func ServePort() string {
	return config.RpcLiPort
}

func init() {
	data, err := cfg.ReadFilePath(configPath)
	if err != nil {
		data = cfg.ReadConfigServer()
		if data == nil {
			log.Fatalln("config.json not exist and configserver was down")
		}
	}
	cfg.Unmarshal(data, &config)

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
