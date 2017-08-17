/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"golang.org/x/net/context"
	"strings"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/dao/mongo"
)

type UserThirdparty struct {
	ID       string "id"
	Name     string "name"
	CloudID  string "cloudid"
	CloudKey string "cloudkey"
}

func (c *CloudTokenService) TokenSet(ctx context.Context, in *pb.CloudTokenCall) (*pb.CloudTokenBack, error) {
	collection := "cloud.token." + strings.ToLower(in.Cloud)
	v := &UserThirdparty{in.AAATokenID, in.Name, in.Id, in.Key}
	ok := mongo.MgoInsert(mongo.MgoCloudConn, v, collection)
	if ok {
		return &pb.CloudTokenBack{Id: in.Id, Key: in.Key}, nil
	}
	return &pb.CloudTokenBack{Id: "", Key: ""}, nil
}

func (c *CloudTokenService) TokenGet(ctx context.Context, in *pb.CloudTokenCall) (*pb.CloudTokenBack, error) {
	collection := "cloud.token." + strings.ToLower(in.Cloud)
	query := mongo.MgoQuerys{"id": in.AAATokenID, "name": in.Name}
	result := query.MgoFindOne(mongo.MgoCloudConn, collection)
	cloudid, ok := result["cloudid"].(string)
	cloudkey, _ := result["cloudkey"].(string)
	if ok {
		return &pb.CloudTokenBack{Id: cloudid, Key: cloudkey}, nil
	}
	return &pb.CloudTokenBack{Id: "", Key: ""}, nil
}
