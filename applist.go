package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/brian1917/vcodeapi"
)

func getApps(credsFile string, allApps bool, txtfile string) []string {
	var apps []string

	if allApps == true {
		appList, err := vcodeapi.ParseAppList(credsFile)
		if err != nil {
			log.Fatal(err)
		}
		for _, app := range appList {
			apps = append(apps, app.AppID)
		}
	} else {
		file, err := os.Open(txtfile)
		if err != nil {
			fmt.Println("error")
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			apps = append(apps, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	return apps
}
