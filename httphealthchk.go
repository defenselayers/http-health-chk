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

func makeURL(urlAddr string, port int, bSSL bool) string {
	var urlPrefix = "http"
	if bSSL {
		urlPrefix = "https"
	}

	if (!bSSL && port != 80) || (bSSL && port != 443) {
		urlAddr = fmt.Sprintf("%s:%d", urlAddr, port)
	}

	url := url.URL{
		Scheme: urlPrefix,
		Host:   urlAddr,
	}
	return url.String()
}

func main() {
	hostPtr := flag.String("host", "127.0.0.1", "URL to check")
	sslPtr := flag.Bool("ssl", false, "enable SSL/TLS connection")
	portPtr := flag.Int("port", 0, "sets TCP port (default 80 for http and 443 for https)")
	timeoutPtr := flag.Int("timeout", 3, "sets timeout in seconds")
	resultPtr := flag.Int("result", 200, "defines expected return code")
	urlpathPtr := flag.String("path", "", "defines additional URL path")
	flag.Parse()

	if *portPtr == 0 {
		if !(*sslPtr) {
			*portPtr = 80
		} else {
			*portPtr = 443
		}
	}

	timeout := time.Duration(time.Duration(*timeoutPtr) * time.Second)
	urlAddr := makeURL(*hostPtr, *portPtr, *sslPtr)

	t := http.Transport{
		Dial: func(network, host string) (net.Conn, error) {
			conn, err := net.DialTimeout(network, host, timeout)
			if err != nil {
				return nil, err
			}
			conn.SetDeadline(time.Now().Add(timeout))
			return conn, nil
		},
	}
	httpClient := http.Client{
		Transport: &t,
	}
	resp, err := httpClient.Get(urlAddr + *urlpathPtr)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	if resp.StatusCode != *resultPtr {
		fmt.Println("Incorrect response code!")
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
