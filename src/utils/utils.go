package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var PHPSESSIDCookie string

func URLToGoQuery(subpath string) (*goquery.Document, func() error, error) {
	req, err := http.NewRequest("GET", "https://www.argyros.com.pa/"+subpath, nil)
	if err != nil {
		return nil, nil, err
	}
	req = AddRequestHeaders(req)
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

func AddRequestHeaders(req *http.Request) *http.Request {
	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", PHPSESSIDCookie))
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "es,en;q=0.9,en-US;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://www.argyros.com.pa")
	req.Header.Set("Referer", "https://www.argyros.com.pa/login.php")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"111\", \"Not(A:Brand\";v=\"8\", \"Chromium\";v=\"111\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Linux\"")

	return req
}
