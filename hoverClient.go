package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/glog"
)

type hoverClient struct {
	username            string
	password            string
	hoverID             string
	hoverToken          *http.Cookie
	hoverTokenTimestamp time.Time
	polltime            int
	hoverIP             string
}

func newHoverClient(config map[string]string) *hoverClient {

	polltime, _ := strconv.Atoi(config["POLLTIME"])

	hc := hoverClient{
		username: config["HOVERUSER"],
		password: config["HOVERPASS"],
		hoverID:  config["HOVERID"],
		polltime: polltime,
	}

	hc.getCurrentHoverIP()

	return &hc
}

func (hc *hoverClient) getAuth() {
	defer glog.Flush()

	glog.Info("attempting authentication with Hover")
	data := map[string]string{"username": hc.username, "password": hc.password}
	datab, _ := json.Marshal(data)

	loginURL := "https://www.hover.com/api/login"
	headers := "application/json"

	response, _ := http.Post(loginURL, headers, bytes.NewReader(datab))

	for _, cookie := range response.Cookies() {

		if cookie.Name == "hoverauth" {
			hc.hoverToken = cookie
			hc.hoverTokenTimestamp = time.Now()
		}

		if hc.hoverToken == nil {
			glog.Fatal("could not authenticate with Hover: ")

		}
	}

	glog.Info("authenticated successfully with Hover")

}

func (hc *hoverClient) checkAuth() {
	defer glog.Flush()

	lastCheckDuration := time.Now().Sub(hc.hoverTokenTimestamp)

	if lastCheckDuration > 6*time.Hour {
		glog.Info("reauthentication needed")
		hc.getAuth()
	}
}

func (hc *hoverClient) getCurrentHoverIP() {
	defer glog.Flush()
	hc.checkAuth()

	glog.Info("checking current IP address setting at Hover")

	url := "https://www.hover.com/api/dns"

	request, _ := http.NewRequest("GET", url, nil)
	request.AddCookie(hc.hoverToken)

	client := &http.Client{}
	resp, _ := client.Do(request)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	// Temp
	// body, _ := ioutil.ReadFile("bodytemp.txt")
	// Temp

	data := []byte(body)

	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		fmt.Println("pooper")
	}

	domains := f.(map[string]interface{})["domains"].([]interface{})

	for _, v := range domains {
		entries := v.(map[string]interface{})["entries"].([]interface{})

		for _, vv := range entries {
			if vv.(map[string]interface{})["id"] == hc.hoverID {
				hc.hoverIP = vv.(map[string]interface{})["content"].(string)
			}
		}
	}
}

// func (hc *hoverClient) call(method, resource, data string) string {
// 	defer glog.Flush()
// 	hc.checkAuth()

// 	url := "https://www.hover.com/api/%s", resource

// 	glog.Info("connecting to Hover")

// 	res, _ := http.Get(url, )

// 	// r = requests.request(method, url, data=data, cookies=self.hoverToken)
// 	// if not r.ok:
// 	// 	logging.error('could not connect to Hover!')
// 	// if r.content:
// 	// 	body = r.json()
// 	// 	if "succeeded" not in body or body["succeeded"] is not True:
// 	// 		logging.error('connected to Hover but data is missing!')
// 	// 	logging.info('request to Hover was successful')
// 	// 	return body

// }
