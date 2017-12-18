package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"strings"
)

type config struct {
	Auth struct {
		CredsFile string `json:"credsFile"`
	} `json:"auth"`

	Scope struct {
		AllApps         bool   `json:"allApps"`
		AppListTextFile string `json:"appListTextFile"`
	} `json:"scope"`

	Mode struct {
		LogOnly          bool `json:"logOnly"`
		ProposeOnly      bool `json:"proposeOnly"`
		ProposeAndAccept bool `json:"proposeAndAccept"`
	} `json:"mode"`

	TargetFlaws struct {
		CWEList           string `json:"cweList"`
		RequireTextInDesc bool   `json:"requireTextInDesc"`
		RequiredText      string `json:"requiredText"`
		Static            bool   `json:"static"`
		Dynamic           bool   `json:"dynamic"`
	} `json:"targetFlaws"`

	MitigationInfo struct {
		MitigationType  string `json:"mitigationType"`
		ProposalComment string `json:"proposalComment"`
		ApprovalComment string `json:"approvalComment"`
	} `json:"mitigationInfo"`
}

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "Veracode username")
}

func parseConfig() config {

	flag.Parse()

	//READ CONFIG FILE
	var config config

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	// CHECK FOR MODE ERRORS
	modeCounter := 0
	if config.Mode.LogOnly == true {
		modeCounter++
	}
	if config.Mode.ProposeOnly == true {
		modeCounter++
	}
	if config.Mode.ProposeAndAccept == true {
		modeCounter++
	}
	if modeCounter > 1 {
		log.Fatal("[!]Only one mode is allowed to be set to true")
	}
	if modeCounter == 0 {
		log.Fatal("[!]At least one mode has to be set to true.")
	}

	// REMOVE SPACES FROM CWE LIST
	if strings.Contains(config.TargetFlaws.CWEList, " ") {
		config.TargetFlaws.CWEList = strings.Replace(config.TargetFlaws.CWEList, " ", "", -1)
	}

	// IF REQUIRED TEXT IS TRUE, CONFIRM TEXT PRESENT
	if config.TargetFlaws.RequireTextInDesc == true && config.TargetFlaws.RequiredText == "" {
		log.Fatal("[!]Need to provide to text to search for in description")
	}

	// CHECK MITIGATION TYPE IS VALID
	if config.MitigationInfo.MitigationType != "appdesign" &&
		config.MitigationInfo.MitigationType != "osenv" &&
		config.MitigationInfo.MitigationType != "netenv" &&
		config.MitigationInfo.MitigationType != "fp" {
		log.Fatal("[!]Mitigation type needs to be appdesign, osenv, netenv, or fp")
	}

	return config
}
