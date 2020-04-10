package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/golang/glog"
)

func main() {

	flag.Set("logtostderr", "true")
	flag.Parse()
	defer glog.Flush()

	config := map[string]string{
		"HOVERUSER": os.Getenv("HOVERUSER"),
		"HOVERPASS": os.Getenv("HOVERPASS"),
		"HOVERID":   os.Getenv("HOVERID"),
		"POLLTIME":  os.Getenv("POLLTIME"),
	}

	for k, v := range config {
		if len(v) == 0 {
			glog.Fatalf("missing environment variable: %s\n", k)
		}
	}

	client := newHoverClient(config)

	for true {
		if client.hoverIP != client.currentIP {
			glog.Infof("hover IP needs to be updated. Hover: %s, Actual: %s",
				client.hoverIP, client.currentIP)

			if err := client.updateHoverIP(client.currentIP); err != nil {
				glog.Error("hover update failed")
			}
		} else {
			glog.Infof("hover IP does not need to be updated. Hover: %s, Actual: %s",
				client.hoverIP, client.currentIP)
		}

		polltime, err := strconv.Atoi(config["POLLTIME"])
		if err != nil {
			polltime = 360
		}

		glog.Infof("sleeping for %d minutes", polltime)
		time.Sleep(time.Duration(polltime) * time.Minute)

	}

}
