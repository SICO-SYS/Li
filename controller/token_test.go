/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"golang.org/x/net/context"
	"testing"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/public"
)

func Test_SetRPC(t *testing.T) {
	test := &CloudTokenService{}
	in := &pb.CloudTokenCall{AAATokenID: "01234567890abcdef", Cloud: "test", Name: "test", Id: "test_accessid", Key: "test_accesskey"}
	res, err := test.SetRPC(context.Background(), in)
	if err != nil {
		t.Error(err)
	}
	if res.Code != 0 {
		t.Error(res.Code)
	}
}

func Benchmark_SetRPC(b *testing.B) {
	test := &CloudTokenService{}
	for i := 0; i < b.N; i++ {
		in := &pb.CloudTokenCall{AAATokenID: public.GenerateHexString(), Cloud: "test", Name: "test", Id: "test_accessid", Key: "test_accesskey"}
		res, err := test.SetRPC(context.Background(), in)
		if err != nil {
			b.Error(err)
		}
		if res.Code != 0 {
			b.Error(res.Code)
		}
	}
}
