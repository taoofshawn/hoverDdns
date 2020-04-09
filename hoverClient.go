package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	currentIP           string
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
	hc.getCurrentExternalIP()

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

type hoverResponse struct {
	Domains []struct {
		Entries []struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		} `json:"entries"`
	} `json:"domains"`
}

func (hc *hoverClient) getCurrentHoverIP() {
	defer glog.Flush()
	hc.checkAuth()

	glog.Info("checking current IP address setting at Hover")

	dnsURL := "https://www.hover.com/api/dns"

	request, _ := http.NewRequest("GET", dnsURL, nil)
	request.AddCookie(hc.hoverToken)

	client := &http.Client{}
	resp, _ := client.Do(request)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	// Temp
	// body, _ := ioutil.ReadFile("bodytemp.txt")
	// ioutil.WriteFile("bodytemp.txt", body, 0644)
	// Temp

	respDump := hoverResponse{}
	json.Unmarshal([]byte(body), &respDump)
	for _, domain := range respDump.Domains {
		for _, entry := range domain.Entries {
			if entry.ID == hc.hoverID {
				hc.hoverIP = entry.Content
			}
		}
	}
}

func (hc *hoverClient) getCurrentExternalIP() {
	defer glog.Flush()
	glog.Info("checking current external IP address")

	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		glog.Error("unable get get current external IP address")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	hc.currentIP = string(body)
}

func (hc *hoverClient) updateHoverIP(newIP string) error {
	defer glog.Flush()
	glog.Infof("updating hover IP with %s", newIP)

	hc.checkAuth()

	dnsURL := "https://www.hover.com/api/dns/" + hc.hoverID
	content := url.Values{}
	content.Add("content", newIP)

	request, _ := http.NewRequest("PUT", dnsURL, strings.NewReader(content.Encode()))
	request.AddCookie(hc.hoverToken)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		glog.Error("could not connect to Hover!")
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if !strings.Contains(string(contents), "\"succeeded\":true") {
		glog.Error("connected to Hover but data is missing!")
	}

	return err

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
