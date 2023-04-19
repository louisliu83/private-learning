package main

import (
	"fmt"
	"strings"

	"io/ioutil"

	"log"

	"net/http"

	"net/url"
)

// 定义所有可用的服务器列表

var proxyList = []string{

	"http://localhost:8091",

	"http://localhost:8092",
}

func main() {

	proxy := &Proxy{

		urls: make([]*url.URL, 0),

		current: 0,

		transport: &http.Transport{},
	}

	for _, p := range proxyList {

		u, _ := url.Parse(p)

		proxy.urls = append(proxy.urls, u)

	}

	// 监听80端口

	if err := http.ListenAndServe(":8080", proxy); err != nil {

		panic(err)

	}

}

type Proxy struct {
	urls []*url.URL

	current int

	transport *http.Transport
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	u := p.urls[0]

	if strings.Contains(r.URL.Path, "/v2/") {
		u = p.urls[1]
	}

	// 设置协议,Authorization和主机

	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

	r.Header.Set("X-Forwarded-Proto", "http")

	r.Header.Set("Authorization", r.Header.Get("Authorization"))

	r.URL.Scheme = u.Scheme

	r.URL.Host = u.Host

	// 发送请求

	res, err := p.transport.RoundTrip(r)

	if err != nil {

		log.Printf("%s", err)

		w.WriteHeader(http.StatusInternalServerError)

		return

	}

	defer res.Body.Close()

	// 将结果返回给客户端

	for k, v := range res.Header {

		w.Header().Set(k, v[0])

	}

	w.WriteHeader(res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {

		log.Printf("%s", err)

		w.WriteHeader(http.StatusInternalServerError)

		return

	}

	fmt.Fprintf(w, "%s", body)

}
