/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"golang.org/x/net/context"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/dao/mongo"
)

type UserThirdparty struct {
	ID       string "id"
	Name     string "name"
	CloudID  string "cloudid"
	CloudKey string "cloudkey"
}

type CloudTokenService struct{}

func (c *CloudTokenService) SetRPC(ctx context.Context, in *pb.CloudTokenCall) (*pb.CloudTokenBack, error) {
	collection := mongo.CollectionCloudTokenName(in.Cloud)
	v := &UserThirdparty{in.AAATokenID, in.Name, in.Id, in.Key}
	err := mongo.Insert(cloudDB, collection, v)
	if err != nil {
		return &pb.CloudTokenBack{Code: 202}, nil
	}
	return &pb.CloudTokenBack{Code: 0}, nil
}

func (c *CloudTokenService) GetRPC(ctx context.Context, in *pb.CloudTokenCall) (*pb.CloudTokenBack, error) {
	collection := mongo.CollectionCloudTokenName(in.Cloud)
	query := mongo.Querys{"id": in.AAATokenID, "name": in.Name}
	result, err := mongo.FindOne(cloudDB, collection, query)
	if err != nil {
		return &pb.CloudTokenBack{Code: 202}, nil
	}
	cloudid, ok := result["cloudid"].(string)
	cloudkey, _ := result["cloudkey"].(string)
	if ok && cloudid != "" {
		return &pb.CloudTokenBack{Code: 0, Id: cloudid, Key: cloudkey}, nil
	}
	return &pb.CloudTokenBack{Code: 2003}, nil
}
