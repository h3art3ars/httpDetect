package main

import (
	"flag"
	"fmt"
	"github.com/h3art3ars/httpDetect/common"
	"os"
)

func main() {
	common.ParseFlag()
	flag.Parse()
	//common.UrlFile = "./url.txt"
	//common.DstPort = "443"
	if common.Ip == "" && common.UrlFile == "" {
		flag.Usage()
		fmt.Println("null host")
		os.Exit(1)
	}
	_, err := common.DetectHttpByHost(common.Ip, common.UrlFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("End")

}
