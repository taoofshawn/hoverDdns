package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"os"
)

func main() {
	
	flag.Set("logtostderr", "true")  // Logging flag
	flag.Parse()
	defer glog.Flush()

	config := map[string]string{
		"HOVERUSER": os.Getenv("HOVERUSER"),
		"HOVERPASS": os.Getenv("HOVERPASS"),
		"HOVERID":   os.Getenv("HOVERID"),
		"POLLTIME":  os.Getenv("POLLTIME"),
	}
	if len(config["POLLTIME"]) == 0 {
		config["POLLTIME"] = "360"
	}

	for k,v := range config {
		if len(v) == 0 {
			glog.Fatalf("missing environment variable: %s\n", k)
		}
	}

	client := newHoverClient(config)
	fmt.Println(client.username)

}
