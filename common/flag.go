package common

import "flag"

func ParseFlag() {
	flag.StringVar(&Ip, "h", "", "url,cidr or single ip eg: 192.168.1.1,192.168.1.0/24")
	flag.StringVar(&UrlFile, "f", "", "url file,cidr or single ip")
	flag.StringVar(&OutputFile, "o", "./res.txt", "output file")
	flag.BoolVar(&SimplePort, "s", false, "simple port scan")
	flag.StringVar(&DstPort, "p", "", "reference port")
	flag.BoolVar(&VerySimplePort, "ss", false, "very simple port scan, 80,443,8000,8080")
	flag.IntVar(&ThreadsAmount, "t", 1500, "goroutine mounts")
	flag.IntVar(&Timeout, "T", 20, "Timeout")
}
