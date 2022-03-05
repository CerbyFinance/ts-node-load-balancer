package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"sync"
	"time"
)

const (
	Attempts int = iota
	Retry
)

type Tokens struct {
	mx     sync.Mutex
	ct     int
	tokens []string
}

func NewTokens(tokens []string) *Tokens {
	return &Tokens{
		tokens: tokens,
		ct:     0,
	}
}

func (c *Tokens) Get() string {
	c.mx.Lock()

	c.ct = (c.ct + 1) % len(c.tokens)
	val := c.tokens[c.ct]

	c.mx.Unlock()
	return val
}

type Proxy struct {
	// URL          *url.URL
	ReverseProxy *httputil.ReverseProxy
	tokens       *Tokens
}

func LinesInFile(fileName string) []string {
	f, _ := os.Open(fileName)
	scanner := bufio.NewScanner(f)
	result := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}
	return result
}

func dialTLS(network, addr string) (net.Conn, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{ServerName: host}

	tlsConn := tls.Client(conn, cfg)
	if err := tlsConn.Handshake(); err != nil {
		conn.Close()
		return nil, err
	}

	cs := tlsConn.ConnectionState()
	cert := cs.PeerCertificates[0]

	// Verify here
	cert.VerifyHostname(host)
	log.Println(cert.Subject)

	return tlsConn, nil
}

var proxies []string = []string{
	"RhbxQP:uEc4dA@104.238.190.248:10296",
	"52qqy7:wX1MNS@85.195.81.143:10122",
	"Uy8j3T:KJWZB2@45.145.57.228:11693",
	"tUEGX8:bRXzV4@45.91.209.140:10484",
	"dh3Ngq:q7BYyD@45.153.20.207:10487",
}

func createProxy(path string) *Proxy {
	serverUrl, _ := url.Parse(path)

	fmt.Println(serverUrl)

	proxy := httputil.NewSingleHostReverseProxy(serverUrl)

	// min := 10001
	// max := 12500
	// port := rand.Intn(max-min) + min

	proxyId := rand.Intn(len(proxies))
	proxyStr := proxies[proxyId]

	// proxyUrl, _ := url.Parse("http://vpsville:Pae9aile@45.139.185.34:" + strconv.Itoa(port))
	proxyUrl, _ := url.Parse(proxyStr)

	// if strings.Contains(path, "eth/mainnet") {
	// 	proxyUrl = nil
	// }

	proxy.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConnsPerHost:   32,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DialTLS:               dialTLS,
	}

	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Host = req.URL.Host
	}

	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
	}

	return &Proxy{
		ReverseProxy: proxy,
		tokens:       NewTokens(tokens),
	}
}

type ProxyMap map[string]*Proxy

func buildBalance(host string) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		fullPath := host + r.URL.Path

		_, ok := proxyMap[fullPath]
		if !ok {
			log.Printf("%s not found, creating\n", fullPath)

			proxyMap[fullPath] = createProxy(host)
		}

		some, _ := proxyMap[fullPath]
		token := some.tokens.Get()

		log.Printf("%s, %s\n", fullPath, token)

		clearedPath := reToken.ReplaceAllString(r.URL.Path, token)

		r.URL, _ = url.Parse(clearedPath)

		some.ReverseProxy.ServeHTTP(w, r)
		return
	}
}

var tokens []string
var proxyMap = make(map[string]*Proxy)
var reToken = regexp.MustCompile(`(\@token)`)

func main() {
	port := 3030

	tokens = LinesInFile("tokens.txt")

	balance := buildBalance("https://speedy-nodes-nyc.moralis.io")

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(balance),
	}

	log.Printf("Balancer started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
