package scrappers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"bitbucket.org/marcoboschetti/catalogscraper/src/utils"
	"github.com/gin-gonic/gin"
)

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
	err = Login(inMemoryUsername, inMemoryPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.Next()
}

func refreshSessionCookie() error {
	if len(utils.PHPSESSIDCookie) > 0 {
		return nil
	}

	resp, err := http.Get("https://www.argyros.com.pa/categories.php")
	if err != nil {
		return err
	}

	cookieHeader := resp.Header.Get("Set-Cookie") //PHPSESSID=7e7938eb5ed884243413d19e4f21e524
	utils.PHPSESSIDCookie = strings.TrimSuffix(strings.TrimPrefix(cookieHeader, "PHPSESSID="), "; path=/")
	return nil
}

func Login(username, password string) error {
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

	req = utils.AddRequestHeaders(req)
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
