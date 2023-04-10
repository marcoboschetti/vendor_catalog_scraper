package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var PHPSESSIDCookie string
var isLogged bool

var inMemoryUsername string
var inMemoryPassword string

func VendorCookieMiddleware(c *gin.Context) {
	err := refreshSessionCookie()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
}

func VendorAuthMiddleware(c *gin.Context) {
	err := refreshSessionCookie()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// Login
	err = login(inMemoryUsername, inMemoryPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.Next()
}

func refreshSessionCookie() error {
	if len(PHPSESSIDCookie) > 0 {
		return nil
	}

	resp, err := http.Get("https://www.argyros.com.pa/categories.php")
	if err != nil {
		return err
	}

	cookieHeader := resp.Header.Get("Set-Cookie") //PHPSESSID=7e7938eb5ed884243413d19e4f21e524
	PHPSESSIDCookie = strings.TrimSuffix(strings.TrimPrefix(cookieHeader, "PHPSESSID="), "; path=/")
	return nil
}

func login(username, password string) error {
	if isLogged {
		return nil
	}

	err := refreshSessionCookie()
	if err != nil {
		return err
	}

	encodedCredentials := url.QueryEscape(fmt.Sprintf("email=%s&password=%s", username, password))
	body := strings.NewReader("usr_login=form_type%3Dcustomer_login%26utf8%3D%25E2%259C%2593%26" + encodedCredentials)

	req, err := http.NewRequest("POST", "https://www.argyros.com.pa/database/data-user.php", body)
	if err != nil {
		return err
	}

	req = addRequestHeaders(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("login error with status code %d: %v", resp.StatusCode, resp)
	}

	// We Read the response body on the line below.
	ansBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if strings.Contains(string(ansBody), "chequee sus credenciales") {
		return errors.New("invalid vendor login credentials")
	}

	isLogged = true
	return nil
}

func addRequestHeaders(req *http.Request) *http.Request {
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
