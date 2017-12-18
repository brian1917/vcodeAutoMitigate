package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
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
		CWEList           bool   `json:"cweList"`
		RequireTextInDesc bool   `json:"requireTextInDesc"`
		RequiredText      string `json:"requiredText"`
		Static            bool   `json:"static"`
		Dynamic           bool   `json:"dynamic"`
	} `json:"targetFlaws"`

	MitigationInfo struct {
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

	return config
}
