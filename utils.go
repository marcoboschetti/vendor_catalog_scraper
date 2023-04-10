package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func urlToGoQuery(subpath string) (*goquery.Document, func() error, error) {
	req, err := http.NewRequest("GET", "https://www.argyros.com.pa/"+subpath, nil)
	if err != nil {
		return nil, nil, err
	}
	req = addRequestHeaders(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("error getting %s with status code %d: %v", subpath, resp.StatusCode, resp)
	}

	// We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, nil, err
	}

	return doc, resp.Body.Close, err
}
