package main

import (
	"httpdetect/common"
)

func main() {
	common.DetectHttp("mohtasib.gov.pk", "443", 30)
}
