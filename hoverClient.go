package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
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
	hoverIP             string
	currentIP           string
}

func newHoverClient(config map[string]string) *hoverClient {


	hc := hoverClient{
		username: config["HOVERUSER"],
		password: config["HOVERPASS"],
		hoverID:  config["HOVERID"],
	}

	hc.getAuth()
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

	response, err := http.Post(loginURL, headers, bytes.NewReader(datab))

	for _, cookie := range response.Cookies() {

		if err != nil || cookie.Name == "hoverauth" {
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

func (hc *hoverClient) call(method, resource, data string) ([]byte, error) {
	defer glog.Flush()
	hc.checkAuth()

	url := "https://www.hover.com/api/" + resource

	glog.Info("connecting to Hover")
	request, _ := http.NewRequest(method, url, strings.NewReader(data))
	request.AddCookie(hc.hoverToken)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		glog.Error("could not connect to Hover!")
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	if !strings.Contains(string(body), "\"succeeded\":true") || err != nil {
		glog.Error("connected to Hover but data is missing!")
	} else {
		glog.Info("request to Hover was successful")
	}

	return body, err
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
	glog.Info("checking current IP address setting at Hover")

	body, _ := hc.call("GET", "dns", "")

	allHoverData := hoverResponse{}
	json.Unmarshal(body, &allHoverData)
	for _, domain := range allHoverData.Domains {
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

	resource := "dns/" + hc.hoverID
	data := "content=" + newIP

	_, err := hc.call("PUT", resource, data)

	return err
}
