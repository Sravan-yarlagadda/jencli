package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Jencli  - Jenkins client
type Jencli struct {
	Crumb struct {
		UsesCrumb   bool
		CrumbString string `json:"crumbRequestField"`
		CrumbValue  string `json:"crumb"`
	}
	User  string
	Token string
}

var tr = &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
}

var client = &http.Client{Transport: tr}

var useCrumbsStruct struct {
	Class     string `json:"_class"`
	UseCrumbs bool   `json:"useCrumbs"`
}

// Start Jenkins job with the given URL
func (cli *Jencli) Start(url string, parameters string, monitor bool) {
	fmt.Println("Start Jenkins job for the given url : ", url)
	jenURLSplit := strings.Split(url, "/")
	jenURL := jenURLSplit[0] + "//" + jenURLSplit[2]
	req, cerr := http.NewRequest("GET", jenURL+"/api/json?tree=useCrumbs", nil)
	if cerr != nil {
		fmt.Println(cerr)
		os.Exit(1)
	}
	req.SetBasicAuth(cli.User, cli.Token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 {
		fmt.Printf("Not Authorized : Error code %d\n", resp.StatusCode)
		os.Exit(1)
	}
	// http://192.168.99.100:8080/api/json?pretty=true&tree=useCrumbs

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Body read error : ", err)
		os.Exit(1)
	}
	fmt.Printf("Response : %s \n", body)
	var useCrumbsVal = useCrumbsStruct
	decErr := json.Unmarshal(body, &useCrumbsVal)
	if decErr != nil {
		fmt.Println("Decode error : ", decErr)
		os.Exit(1)
	}
	fmt.Printf("class : %s\nuseCrumbs : %t\n", useCrumbsVal.Class, useCrumbsVal.UseCrumbs)
	// using crumbs. Populate crumb values in struct
	if useCrumbsVal.UseCrumbs == true {
		getCrumbReq, getCrumbReqErr := http.NewRequest("GET", jenURL+"/crumbIssuer/api/json", nil)
		if getCrumbReqErr != nil {
			fmt.Println(getCrumbReqErr)
			os.Exit(1)
		}
		getCrumbReq.SetBasicAuth(cli.User, cli.Token)
		getCrumbResp, getCrumbRespErr := client.Do(getCrumbReq)
		if getCrumbRespErr != nil {
			fmt.Println(getCrumbRespErr)
			os.Exit(1)
		}
		if getCrumbResp.StatusCode != 200 {
			fmt.Println("Http Error : ", getCrumbResp.StatusCode)
			os.Exit(1)
		}
		defer getCrumbResp.Body.Close()
		getCrumbRespBody, _ := ioutil.ReadAll(getCrumbResp.Body)
		// fmt.Printf("getCrumbResp Body : %s\n", getCrumbRespBody)
		crumbDecErr := json.Unmarshal(getCrumbRespBody, &cli.Crumb)
		if crumbDecErr != nil {
			fmt.Println(crumbDecErr)
			os.Exit(1)
		}
		cli.Crumb.UsesCrumb = true
	}

	startJobReq, startJobReqErr := http.NewRequest("POST", url+"/build", nil)
	startJobReq.SetBasicAuth(cli.User, cli.Token)
	if startJobReqErr != nil {
		log.Fatalln(startJobReqErr)
	}
	if cli.Crumb.UsesCrumb == true {
		// fmt.Printf("setting curmb header : %s : %s\n", cli.Crumb.CrumbString, cli.Crumb.CrumbValue)
		startJobReq.Header.Set(cli.Crumb.CrumbString, cli.Crumb.CrumbValue)

	}
	nextBuildNum := cli.getNextBuildNumber(url)
	startJobResp, startJobRespErr := client.Do(startJobReq)
	if startJobRespErr != nil {
		log.Fatalln(startJobRespErr)
	}
	if startJobResp.StatusCode != 201 {
		log.Fatalln("Unable to start Job. Got status code : ", startJobResp.StatusCode)
	}
	log.Printf("Job %s Queued..\n", url)

	if monitor == true {
		isBuildStarted := false
		for isBuildStarted == false {
			newNextBuildNum := cli.getNextBuildNumber(url)
			if newNextBuildNum > nextBuildNum {
				isBuildStarted = true
			}
			time.Sleep(10 * time.Second)
		}
		log.Printf("Build #%v Started for job %s\n", nextBuildNum, url)
	}

}

func (cli *Jencli) monitorBuild(url string, buildNumber int) {

}

func (cli *Jencli) getNextBuildNumber(url string) int {
	req, _ := http.NewRequest("GET", url+"/api/json?tree=nextBuildNumber", nil)
	req.SetBasicAuth(cli.User, cli.Token)
	if cli.Crumb.UsesCrumb == true {
		req.Header.Set(cli.Crumb.CrumbString, cli.Crumb.CrumbValue)
	}
	res, respErr := client.Do(req)
	if respErr != nil {
		log.Fatalln(respErr)
	}

	type nextBuildNumber struct {
		Class           string `json:"_class"`
		NextBuildNumber int    `json:"nextBuildNumber"`
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var nextBuild = nextBuildNumber{}
	json.Unmarshal(body, &nextBuild)
	//fmt.Println(nextBuild.NextBuildNumber)
	return nextBuild.NextBuildNumber
}
