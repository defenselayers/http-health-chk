/* httphealtcheck (c) Aleksander P. Czarnowski, Defenselayers Sp. z o.o. 2020
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

var timeout = time.Duration(time.Second)

func Timeout(network, host string) (net.Conn, error) {
	conn, err := net.DialTimeout(network, host, timeout)
	if err != nil {
		return nil, err
	}
	conn.SetDeadline(time.Now().Add(timeout))
	return conn, nil
}

func makeURL(urlAddr string, Port int, bSSL bool) string {
	var urlPrefix = "http"
	if bSSL {
		urlPrefix = "https"
	}

	if Port != 80 {
		urlAddr = fmt.Sprintf("%s:%d", urlAddr, Port)
	}

	url := url.URL{
		Scheme: urlPrefix,
		Host:   urlAddr,
	}
	return url.String()
}

func main() {
	var exitCode int = 0
	var urlAddr string

	hostPtr := flag.String("host", "127.0.0.1", "URL to check")
	sslPtr := flag.Bool("ssl", false, "enable SSL/TLS connection")
	portPtr := flag.Int("port", 80, "sets TCP port")
	timeoutPtr := flag.Int("timeout", 3, "sets timeout in seconds")
	resultPtr := flag.Int("result", 200, "defines expected return code")
	urlpathPtr := flag.String("path", "", "defines additional URL path")
	flag.Parse()

	timeout = time.Duration(time.Duration(*timeoutPtr) * time.Second)
	urlAddr = makeURL(*hostPtr, *portPtr, *sslPtr)

	t := http.Transport{
		Dial: Timeout,
	}
	httpClient := http.Client{
		Transport: &t,
	}
	//resp, err := http.Get(urlAddr)
	resp, err := httpClient.Get(urlAddr + *urlpathPtr)
	if err != nil {
		// log.Panicln(err)
		log.Fatal(err)
	}
	defer resp.Body.Close()
	//if resp.StatusCode != 200 {
	//	b, _ := ioutil.ReadAll(resp.Body)
	//	log.Fatal(string(b))
	//}
	fmt.Println(resp.StatusCode)

	if resp.StatusCode != *resultPtr {
		fmt.Println("Wrong response code!")
		exitCode = 1
	}

	os.Exit(exitCode)
}
