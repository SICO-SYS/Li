/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package main

import (
	"net"

	"github.com/SiCo-Ops/Li/controller"
	"github.com/SiCo-Ops/cfg"
)

func Run() {
	lis, _ := net.Listen("tcp", cfg.Config.Rpc.Li)
	controller.S.Serve(lis)
}

func main() {
	Run()
}
