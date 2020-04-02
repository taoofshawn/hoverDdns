package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/glog"
)

type hoverClient struct {
	username            string
	password            string
	hoverID             string
	hoverToken          string
	hoverTokenTimestamp time.Time
	polltime            int
}

func newHoverClient(config map[string]string) *hoverClient {

	polltime, _ := strconv.Atoi(config["POLLTIME"])

	hc := hoverClient{
		username:   config["HOVERUSER"],
		password:   config["HOVERPASS"],
		hoverID:    config["HOVERID"],
		hoverToken: "",
		polltime:   polltime,
	}

	return &hc
}

func (hc *hoverClient) getAuth() {
	glog.Info("attempting authentication with Hover")
	data := map[string]string{"username": hc.username, "password": hc.password}
	datab, _ := json.Marshal(data)

	loginURL := "https://www.hover.com/api/login"
	headers := "application/json"

	res, _ := http.Post(loginURL, headers, bytes.NewReader(datab))

	for _, cookie := range res.Cookies() {
		if cookie.Name == "hoverauth" {
			hc.hoverToken = cookie.Value
			hc.hoverTokenTimestamp = time.Now()
		}
	}
}

