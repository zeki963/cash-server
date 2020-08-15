package notify

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

//Postnotifyslack 傳訊息到slack
func Postnotifyslack(text string) {
	authURL := "https://hooks.slack.com/services/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"
	toServerVal := "payload={\"text\": \"[ " + getname() + getip() + " ] > " + text + "\"}"
	resp, err := http.Post(authURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(toServerVal))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
}

func getip() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "ip"
}

func getname() string {
	name, err := os.Hostname()
	if err != nil {
		return "name / "
	}

	return name + " / "
}
