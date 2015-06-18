package token

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func GetTokenFromUAA(username string, password string, uaaURL string, skipSSLVerification bool) (string, error) {
	config := &tls.Config{InsecureSkipVerify: skipSSLVerification}
	data := url.Values{
		"grant_type": {"password"},
		"username":   {username},
		"password":   {password},
		"client_id":  {"cf"},
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/oauth/token", uaaURL), strings.NewReader(data.Encode()))

	request.SetBasicAuth("cf", "")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tr := &http.Transport{TLSClientConfig: config}
	httpClient := &http.Client{Transport: tr}

	response, err := httpClient.Do(request)
	if err != nil {
		return "", err
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	response.Body.Close()

	uaaResponse := make(map[string]interface{})
	err = json.Unmarshal(responseBytes, &uaaResponse)
	if err != nil {
		return "", err
	}

	authToken := uaaResponse["token_type"].(string) + " " + uaaResponse["access_token"].(string)

	return authToken, nil
}
