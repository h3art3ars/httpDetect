package main

import (
	"flag"
	"fmt"
	"github.com/h3art3ars/httpDetect/common"
)

func main() {
	common.ParseFlag()
	flag.Parse()
	//common.Ip = "149.54.1.220"
	//common.DstPort = "443"
	_, err := common.DetectHttpByHost(common.Ip, common.UrlFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("End")

}
