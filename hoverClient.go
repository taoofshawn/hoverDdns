package main

import (
	"strconv"
)

type hoverClient struct {
	username string
	password string
	hoverid  string
	polltime int
}

func newHoverClient(config map[string]string) *hoverClient {

	polltime, _ := strconv.Atoi(config["POLLTIME"])

	hc := hoverClient{
		username: config["HOVERUSER"],
		password: config["HOVERPASS"],
		hoverid:  config["HOVERID"],
		polltime: polltime,
	}

	return &hc
}

func getAuth(hc *hoverClient) {

}
