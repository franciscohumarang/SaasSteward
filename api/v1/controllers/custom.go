package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ExternalEndpoint(c *gin.Context, url string) {

	// create a new HTTP client
	client := &http.Client{}

	// create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating HTTP request")
		return
	}

	// send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error sending HTTP request")
		return
	}
	defer resp.Body.Close()

	// read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading response body")
		return
	}

	// set the response status code and body
	//c.String(resp.StatusCode, string(body))
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, body, "", "\t")
	if error != nil {

		return
	}

}
