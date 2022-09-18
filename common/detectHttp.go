package common

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var ports = strings.Split(WebPort, ",")
var (
	routineAmount chan struct{}
	urlPool       chan string
	urlValid      []string
	wg            sync.WaitGroup
	File          *os.File
)

func DetectHttpByHost(host string, filename string) ([]string, error) {
	routineAmount = make(chan struct{}, ThreadsAmount)
	urlPool = make(chan string, len(ports)*len(host))
	hosts, err := ParseIP(host, filename)
	filePath := OutputFile
	File, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer func() {
		close(urlPool)
		File.Close()
	}()
	if err != nil {
		return nil, err
	}
	go func() {
		for url := range urlPool {

			if url != "" {
				r := bufio.NewWriter(File)
				_, err := r.Write([]byte(url + "\n"))
				if err != nil {
					fmt.Println(err)
				}
				r.Flush()
				fmt.Printf("[+]URL:\t%s\n", url)
				urlValid = append(urlValid, url)
			}
			<-routineAmount
			wg.Done()
		}
	}()
	if SimplePort {
		ports = strings.Split(WebPortSimple, ",")
	} else if VerySimplePort {
		ports = strings.Split("80,443,8000,8080", ",")
	}
	if DstPort != "" {
		if strings.Contains(DstPort, ",") {
			ports = strings.Split(DstPort, ",")
		} else {
			ports = []string{DstPort}
		}

	}
	if len(hosts) < 1 {
		return nil, errors.New("null host")
	}
	wg.Add(len(ports) * len(hosts))
	for _, p := range ports {
		for _, h := range hosts {
			routineAmount <- struct{}{}
			go detect(h, p, urlPool)
		}
	}
	wg.Wait()
	return urlValid, nil
}

//阻塞方式检测port
func detect(host, port string, res chan string) {
	schema, err := DetectHttp(host, port, Timeout)
	if err != nil {
		//fmt.Println(err)
		res <- ""
		return
	}
	res <- fmt.Sprintf("%s://%s:%s", schema, host, port)
	return
}
func DetectHttp(host string, port string, duration int) (string, error) {
	//fmt.Printf("%s:%s\n", host, port)
	//var duration time.Duration
	if duration == 0 {
		duration = 5
	}

	ipPort := fmt.Sprintf("%s:%s", host, port)
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    time.Duration(duration) * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(duration) * time.Second,
	}
	resp, err := client.Get(fmt.Sprintf("https://%s", ipPort))
	if err != nil || resp.TLS == nil {
		goto detectHttp
	}
	defer resp.Body.Close()
	return "https", nil
	//tcp connect

detectHttp:
	client = &http.Client{
		Timeout: time.Duration(duration) * time.Second,
		//CheckRedirect: func(req *http.Request, via []*http.Request) error {
		//	fmt.Printf("")
		//	return nil
		//},
	}
	resp, err = client.Get(fmt.Sprintf("http://%s", ipPort))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return "http", err

}
func DetectHttpBak(host string, port string, duration int) (string, error) {
	//fmt.Printf("%s:%s\n", host, port)
	//var duration time.Duration
	if duration == 0 {
		duration = 5
	}

	ipPort := fmt.Sprintf("%s:%s", host, port)
	httpSender := "GET / HTTP/1.1\r\n"
	//tcp connect
	tcpAddr, err := net.ResolveTCPAddr("tcp", ipPort)
	if err != nil {
		return "", err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	//detect http
	_, err = conn.Write([]byte(httpSender))
	if err != nil {
		return "", err
	}
	//超时

	err = conn.SetDeadline(time.Now().Add(time.Duration(Timeout) * time.Second))
	if err != nil {
		return "", err
	}
	r := make([]byte, 100)
	_, err = conn.Read(r)
	if err != nil {
		conn.Close()
		goto detectHttps
		//return "", err
	}
	if bytes.HasPrefix(r, []byte("HTTP")) {
		return "http", nil
	}
detectHttps:
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn2, err := tls.Dial("tcp", ipPort, tlsConfig)
	if err != nil {
		return "", err
	}
	_, err = conn2.Write([]byte(httpSender))
	if err != nil {
		return "", err
	}
	defer conn2.Close()
	err = conn2.SetDeadline(time.Now().Add(time.Duration(Timeout) * time.Second))
	if err != nil {
		return "", err
	}
	r1 := make([]byte, 100)
	_, err = conn2.Read(r1)
	if err != nil {
		return "", err
	}
	if bytes.HasPrefix(r1, []byte("HTTP")) {
		return "https", nil
	}
	return "", errors.New("unknown protocol")

}
