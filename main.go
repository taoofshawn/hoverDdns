package main

import (
	"fmt"
	"os"
)

func main() {
	config := map[string]string{
		"HOVERUSER": os.Getenv("HOVERUSER"),
		"HOVERPASS": os.Getenv("HOVERPASS"),
		"HOVERID": os.Getenv("HOVERID"),
		"POLLTIME": os.Getenv("POLLTIME"),
		"LOGLEVEL": os.Getenv("LOGLEVEL"),
	}
	if len(config["POLLTIME"]) == 0 {
		config["POLLTIME"] = "360"
	}
	if len(config["LOGLEVEL"]) == 0 {
		config["LOGLEVEL"] = "INFO"
	}

	fmt.Println("map: ", config)

}

