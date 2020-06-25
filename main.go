package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/utilitywarehouse/sys-graylog-configurer/graylog"
)

type Client struct {
	username   string
	password   string
	baseURL    string
	httpClient *http.Client
}

func main() {
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	admins := os.Getenv("ADMINS")
	if adminPassword == "" {
		fmt.Println("ADMIN_PASSWORD not set")
		os.Exit(1)
	}

	cl := graylog.NewClient("http://127.0.0.1:9000/api", adminPassword)

	waitForAPI(cl)

	setAdmins(cl, admins)
	fmt.Println("Going to sleep")
}

func setAdmins(cl *graylog.Client, admins string) error {
	if admins == "" {
		return nil
	}
	names := strings.Split(admins, " ")
	for _, n := range names {
		err := cl.SetAdmin(n)
		if err != nil {
			fmt.Println("Error:", err)
		}

	}
	return nil
}

func waitForAPI(cl *graylog.Client) {
	for {
		err := cl.ApiReachable()
		if err == nil {
			break
		}
		fmt.Printf("Graylog API not ready, sleeping for 3 seconds...(%s)\n", err)
		time.Sleep(3 * time.Second)
	}
}
