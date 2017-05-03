/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	"google.golang.org/grpc"
	"io/ioutil"
	// "net/http"
	// "net/url"

	"github.com/SiCo-DevOps/Pb"
	. "github.com/SiCo-DevOps/log"
	// "github.com/SiCo-DevOps/public"
)

var (
	S             = grpc.NewServer()
	assert_Action map[string]interface{}
	// nonce = public.GenNonce()
	// ts    = public.TS()
)

type Cloud_API struct{}

func transAction(cloud string, bsns string, action string) (string, bool) {
	d, err := ioutil.ReadFile("action.json")
	if err != nil {
		LogFatalMsg(0, "controller.transAction")
	}
	json.Unmarshal(d, &assert_Action)

	get_cloud, ok := assert_Action[action].(map[string]interface{})
	if ok {
		get_bsns, ok := get_cloud[cloud].(map[string]interface{})
		if ok {
			value, ok := get_bsns[bsns].(string)
			if ok {
				return value, true
			}
			return "error bsns", false
		}
		return "error cloud", false
	}
	return "error action", false
}

// func Cloud_POST(urlstr string, data *url.Values) {
// 	http.PostForm(urlstr, data)

// }

func init() {
	pb.RegisterCloud_APIServer(S, &Cloud_API{})
}
