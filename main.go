package main

import (
	"flag"
	"fmt"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var allowOnceList = make([]string, 0)

func main() {

	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8080", "proxy listen address")
	flag.Parse()
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().DoFunc(handleRequest)
	proxy.OnResponse().DoFunc(handleResponse)
	proxy.Verbose = *verbose
	log.Fatal(http.ListenAndServe(*addr, proxy))
}

func handleRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

	log.Printf("Request: %s\n", req.URL)
	referer := req.Header.Get("Referer")
	fmt.Println("referer: " + referer)
	// allow normal website loading
	if referer == "" {
		log.Print("allow (empty referer)")
		return req, nil
	}

	refererUrl, err := url.Parse(referer)
	if err != nil {
		fmt.Printf("failed to parse url %e\n", err)
	}
	// allow local link
	if req.Host == refererUrl.Host {
		log.Printf("allowing (local link) from %s to %s\n", referer, req.Host)
		return req, nil
	}
	index := find(allowOnceList, req.URL.String())
	if index != -1 {
		log.Printf("allowing %s once", req.URL.String())
		allowOnceList = append(allowOnceList[:index], allowOnceList[index+1:]...)
		return req, nil
	}
	log.Printf("blocked")
	return req, goproxy.NewResponse(req, "text/html", 200, "Blocked request")
}

func handleResponse(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if resp == nil {
		return resp
	}

	contentType := resp.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "text") {

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read body: %e\n", err)
			return resp
		}
		body := string(data)
		extractLinks(body)

		newResp := goproxy.NewResponse(resp.Request, contentType, resp.StatusCode, body)
		return newResp
	}
	return resp
}

// search for external links and add them to the allowOnceList
func extractLinks(body string) {

	srcLink := regexp.MustCompile(`src="([^"]*)"`)
	matches := srcLink.FindAllStringSubmatch(body, 100)

	for _, match := range matches {
		if len(match) > 1 {
			link := match[1]

			if isExternal(link) && find(allowOnceList, link) == -1 {
				println("appending " + link)
				allowOnceList = append(allowOnceList, link)
			}
		}
	}

	fmt.Println("allowOnceList:")
	for _, item := range allowOnceList {
		fmt.Println("\t" + item)
	}
}


func isExternal(link string) bool {
	// extremely stupid way to test this
	return strings.HasPrefix(link, "http") || strings.HasPrefix(link, "www")
}

func find(array []string, item string) int {
	for i, current := range array {
		if current == item {
			return i
		}
	}
	return -1
}