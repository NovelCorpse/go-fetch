package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/net/html"
	"crypto/tls"
)

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

func fetch() (out []byte, e error) {
	for _, url := range os.Args[1:] {
		// skip tls verify;
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		if resp, e := http.Get(url); e != nil {
			return nil, e
		} else {
			b, e := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if e != nil {
				return nil, e
			}
			out = append(out, b...)

			return out, nil
		}
	}
	return 
}

func main() {
	data, e := fetch()
	if e != nil {
		fmt.Fprintf(os.Stderr, "error with fetch data from url: %s\n", e)
		os.Exit(1)
	}

	r := bytes.NewReader(data)

	doc, e := html.Parse(r)
	if e != nil {
		fmt.Fprintf(os.Stderr, "link-fetcher: %v\n", e)
		os.Exit(1)
	}

	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}

}
