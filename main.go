/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package main

import (
	"net"

	"github.com/SiCo-DevOps/Li/controller"
	"github.com/SiCo-DevOps/cfg"
)

func Run() {
	lis, _ := net.Listen("tcp", cfg.Config.Rpc.Li)
	controller.S.Serve(lis)
}

func main() {
	Run()
}
