package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
//func (cli *Jencli) Start ( url string, parameters string, monitor bool)
//parameters : 	l - url to start the job
//				parameters - parameters JSON format to be passed to jenkins job
//				monitor - boolean value to specify to monitor the build
func (cli *Jencli) Start(l string, parameters string, monitor bool) {
	fmt.Println("Start Jenkins job for the given url : ", l)
	jenURLSplit := strings.Split(l, "/")
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Body read error : ", err)
		os.Exit(1)
	}
	var useCrumbsVal = useCrumbsStruct
	decErr := json.Unmarshal(body, &useCrumbsVal)
	if decErr != nil {
		fmt.Println("Decode error : ", decErr)
		os.Exit(1)
	}
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
		crumbDecErr := json.Unmarshal(getCrumbRespBody, &cli.Crumb)
		if crumbDecErr != nil {
			fmt.Println(crumbDecErr)
			os.Exit(1)
		}
		cli.Crumb.UsesCrumb = true
	}
	var startJobReq *http.Request
	var startJobReqErr error
	if strings.TrimSpace(parameters) != "" {
		form := url.Values{}
		form.Add("json", parameters)
		startJobReq, startJobReqErr = http.NewRequest("POST", l+"/build", strings.NewReader(form.Encode()))
		startJobReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		startJobReq, startJobReqErr = http.NewRequest("POST", l+"/build", nil)
	}
	startJobReq.SetBasicAuth(cli.User, cli.Token)
	if startJobReqErr != nil {
		log.Fatalln(startJobReqErr)
	}
	if cli.Crumb.UsesCrumb == true {
		startJobReq.Header.Set(cli.Crumb.CrumbString, cli.Crumb.CrumbValue)

	}
	nextBuildNum := cli.getNextBuildNumber(l)

	fmt.Printf("Request to start job : %v\n", startJobReq.Form)
	startJobResp, startJobRespErr := client.Do(startJobReq)
	if startJobRespErr != nil {
		log.Fatalln(startJobRespErr)
	}
	if startJobResp.StatusCode != 201 {
		body, _ := ioutil.ReadAll(startJobResp.Body)
		log.Printf("Body : \n %s \n", body)
		log.Fatalln("Unable to start Job. Got status code : ", startJobResp.StatusCode)

	}
	log.Printf("Job %s Queued..\n", l)

	isBuildStarted := false
	for isBuildStarted == false {
		newNextBuildNum := cli.getNextBuildNumber(l)
		if newNextBuildNum > nextBuildNum {
			isBuildStarted = true
		}
		time.Sleep(10 * time.Second)
	}
	log.Printf("Build #%v Started for job %s\n", nextBuildNum, l)
	if monitor == true {
		cli.monitorBuild(l, nextBuildNum)
	}

}

func (cli *Jencli) monitorBuild(url string, buildNumber int) {
	log.Printf("Monitoring build %v for job %v\n", buildNumber, url)
	fmt.Println("\t***************************************************")
	fmt.Printf("\tbuild log : %s\n", url+"/"+strconv.Itoa(buildNumber))
	fmt.Println("\t***************************************************")
	pos := 0
	hasMoreData := true
	for hasMoreData == true {
		time.Sleep(5 * time.Second)
		req, err := http.NewRequest("GET", url+"/"+strconv.Itoa(buildNumber)+"/logText/progressiveText?start="+strconv.Itoa(pos), nil)
		if err != nil {
			log.Fatalln(err)
		}
		req.SetBasicAuth(cli.User, cli.Token)
		if cli.Crumb.UsesCrumb == true {
			req.Header.Set(cli.Crumb.CrumbString, cli.Crumb.CrumbValue)
		}
		res, resErr := client.Do(req)
		if resErr != nil {
			log.Fatalln(resErr)
		}
		if res.StatusCode != 200 {
			log.Fatalln("Error tailing log : ", res.StatusCode)
		}
		defer res.Body.Close()
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatalln(readErr)
		}
		fmt.Printf("%s", body)
		hasMoreData, _ = strconv.ParseBool(res.Header.Get("X-More-Data"))
		pos, _ = strconv.Atoi(res.Header.Get("X-Text-Size"))
	}
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
	return nextBuild.NextBuildNumber
}
