package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//	"fmt"
	"net/http"
	"strconv"
	"time"
	// "reflect"

	"github.com/golang/glog"
)

type hoverClient struct {
	username            string
	password            string
	hoverID             string
	hoverToken          *http.Cookie
	hoverTokenTimestamp time.Time
	polltime            int
}

func newHoverClient(config map[string]string) *hoverClient {

	polltime, _ := strconv.Atoi(config["POLLTIME"])

	hc := hoverClient{
		username: config["HOVERUSER"],
		password: config["HOVERPASS"],
		hoverID:  config["HOVERID"],
		polltime: polltime,
	}

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

	b := []byte(body)

	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		fmt.Println("pooper")
	}

	m := f.(map[string]interface{})

	for k, v := range m {
		if k == "domains" {
			for _, vv := range v.([]interface{}) {
				uu := vv.(map[string]interface{})
				entries := uu["entries"]
				for _, e := range entries.([]interface{}) {
					entry := e.(map[string]interface{})
					if entry["id"] == hc.hoverID {
						fmt.Println(entry["content"]) //Yikes
					}
				}
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
