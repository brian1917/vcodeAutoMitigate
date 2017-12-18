package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/brian1917/vcodeapi"
)

func main() {

	// SET UP LOGGING FILE
	f, err := os.OpenFile("vcodeAutoMitigate"+time.Now().Format("20060102_150405")+".log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Printf("Started running")

	// SET SOME VARIABLES
	var appSkip bool
	var flaws []vcodeapi.Flaw
	var recentBuild string
	var errorCheck error
	var flawList []string
	var buildsBack int

	// PARSE CONFIG FILE AND LOG CONFIG SETTINGS
	config := parseConfig()
	log.Printf("Config Settings:\n"+
		"[*] Creds File: %v\n"+
		"[*] Scope:%v\n"+
		"[*] Mode Text:%v\n"+
		"[*] Target Flaws:%v\n"+
		"[*] Mitigation Info:%v\n",
		config.Auth.CredsFile, config.Scope, config.Mode, config.TargetFlaws, config.MitigationInfo)

	// GET APP LIST
	appList := getApps(config.Auth.CredsFile, config.Scope.AllApps, config.Scope.AppListTextFile)
	appCounter := 0

	// CYCLE THROUGH EACH APP
	for _, appID := range appList {
		//ADJUST SOME VARIABLES
		flawList = []string{}
		appSkip = false
		appCounter++

		fmt.Printf("Processing App ID %v (%v of %v)\n", appID, appCounter, len(appList))

		//GET THE BUILD LIST
		buildList, err := vcodeapi.ParseBuildList(config.Auth.CredsFile, appID)
		if err != nil {
			log.Fatal(err)
		}

		// GET FOUR MOST RECENT BUILD IDS
		if len(buildList) == 0 {
			appSkip = true
			flaws = nil
			recentBuild = ""
		} else {
			//GET THE DETAILED RESULTS FOR MOST RECENT BUILD
			flaws, _, errorCheck = vcodeapi.ParseDetailedReport(config.Auth.CredsFile, buildList[len(buildList)-1].BuildID)
			recentBuild = buildList[len(buildList)-1].BuildID
			buildsBack = 1
			//IF THAT BUILD HAS AN ERROR, GET THE NEXT MOST RECENT (CONTINUE FOR 4 TOTAL BUILDS)
			for i := 1; i < 4; i++ {
				if len(buildList) > i && errorCheck != nil {
					flaws, _, errorCheck = vcodeapi.ParseDetailedReport(config.Auth.CredsFile, buildList[len(buildList)-(i+1)].BuildID)
					recentBuild = buildList[len(buildList)-(i+1)].BuildID
					buildsBack = i + 1
					fmt.Println(buildsBack)
				}
			}
			// IF 4 MOST RECENT BUILDS HAVE ERRORS, THERE ARE NO RESULTS AVAILABLE
			if errorCheck != nil {
				appSkip = true
			}
		}

		//CHECK FLAWS AND
		if appSkip == false {
			for _, f := range flaws {
				// ONLY RUN ON OPEN, RE-OPENED, AND NEW FLAWS (CONFIRM!!!)
				if f.RemediationStatus == "Open" || f.RemediationStatus == "Re-opened" || f.RemediationStatus == "New" {
					// CHECK FOR INCLUDING COMMENT TEXT
					if config.TargetFlaws.RequireTextInDesc == true && strings.Contains(f.Description, config.TargetFlaws.RequiredText) {
						// CHECK IF EXPIRED AND BUILD ARRAY
						flawList = append(flawList, f.Issueid)
					}
				}

			}
			// IF WE HAVE FLAWS MEETING CRITERIA, RUN UPDATE MITIGATION API
			if len(flawList) > 0 {
				log.Printf("[*]Trying to mitigate IDs %v in Build ID %v in App ID %v", flawList, recentBuild, appID)

				if config.Mode.LogOnly == true {
					log.Printf("App ID: %v: Flaw ID(s) %v meet criteria\n", appID, strings.Join(flawList, ","))
				} else {

					// PROPOSE AND MITIGATE THE FLAWS
					actions := [2]string{"propose", "approve"}
					for i := 0; i <= 1; i++ {
						expireError := vcodeapi.ParseUpdateMitigation(config.Auth.CredsFile, recentBuild,
							actions[i], config.MitigationInfo.ProposalComment, strings.Join(flawList, ","))
						// IF WE HAVE AN ERROR, WE NEED TO TRY 2 BUILDS BACK FROM RESULTS BUILD
						// EXAMPLE = RESULTS IN BUILD 3 (MANUAL); DYNAMIC IS BUILD 2; STATIC IS BUILD 1 (BUILD WE NEED TO MITIGATE STATIC FLAW)
						for i := 0; i < 1; i++ {
							if expireError != nil {
								expireError = vcodeapi.ParseUpdateMitigation(config.Auth.CredsFile, recentBuild,
									actions[i], config.MitigationInfo.ProposalComment, strings.Join(flawList, ","))
								if expireError != nil {
									log.Printf("[*] %v", expireError)
								}
							}
						}
						// IF EXPIRE ERROR IS STILL NOT NULL, NOW WE LOG THE ERROR AND EXIT
						if expireError != nil {
							log.Fatalf("[!] Could not "+actions[i]+" mitigation for Flaw IDs %v in App ID %v", flawList, appID)
						}
						// LOG SUCCESSFUL PROPOSED MITIGATIONS
						log.Printf("App ID %v: "+actions[i]+" Flaw IDs %v\n", appID, strings.Join(flawList, ","))
					}
				}
			}
		}
	}
	log.Printf("Completed running")
}
