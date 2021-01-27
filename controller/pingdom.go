package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	extensions "k8s.io/api/extensions/v1beta1"
	utils "pingdom_controller/general"
)

type pingdomClient struct{
	endpoint string
	token string
}

type pingdomCheck struct {
	Name           string   `json:"name"`
	Host           string   `json:"host"`
	Integrationids []int    `json:"integrationids"`
	ProbeFilters   []string `json:"probe_filters"`
	CustomPath     string   `json:"url"`
}

type newPingdomCheck struct {
	pingdomCheck
	Type           string   `json:"type"`
}

type responsePingdomCheck struct {
	Checks []struct {
	ID                int    `json:"id"`
	Created           int    `json:"created"`
	Name              string `json:"name"`
	Hostname          string `json:"hostname"`
	Resolution        int    `json:"resolution"`
	Type              string `json:"type"`
	Ipv6              bool   `json:"ipv6"`
	VerifyCertificate bool   `json:"verify_certificate"`
	Lasterrortime     int    `json:"lasterrortime"`
	Lasttesttime      int    `json:"lasttesttime"`
	Lastresponsetime  int    `json:"lastresponsetime"`
	Status            string `json:"status"`
	} `json:"checks"`
	Counts struct {
	Total    int `json:"total"`
	Limited  int `json:"limited"`
	Filtered int `json:"filtered"`
	} `json:"counts"`
}

type PingdomEngine struct {
	incomingIngress  chan *extensions.Ingress
	deleteIngress  chan string
}

func NewPingdomEngine() *PingdomEngine {
	return &PingdomEngine{
		incomingIngress:  make(chan *extensions.Ingress),
		deleteIngress: make(chan string),
	}
}
 var pClient = &pingdomClient{
 	 endpoint: os.Getenv("PINGDOM_ENDPOINT"),
	 token: os.Getenv("PINGDOM_TOKEN"),
 }

func (p *PingdomEngine) Run() {
	for {
		select {
		case ing := <-p.incomingIngress:
			applyNewCheck(ing)
		case ingName := <-p.deleteIngress:
			deleteCheck(ingName)
		}
	}
}

func applyNewCheck(ing *extensions.Ingress) {
	integrationidsAnnotate, err := utils.GetAnnotationValue(ing.Annotations, "integrationids")
	var integrationids = []int{}
	if (err != nil) || (integrationidsAnnotate == ""){
		log.Printf("No integrationids input")
	}else{
		integrationidsSplit := strings.Split(integrationidsAnnotate, ",")
		for _, integration := range integrationidsSplit {
			res, err := strconv.Atoi(integration)
			if err != nil{
				log.Fatal("Cannot convert to int in Integrationids")
				break
			}
			integrationids = append(integrationids, res)
		}
	}

	probeFiltersAnnotate, err := utils.GetAnnotationValue(ing.Annotations, "probe-filters")
	var probeFilters = []string{}
	if err != nil {
		probeFilters = append(probeFilters, "region: " + os.Getenv("PINGDOM_PROBE_FILTERS"))
	} else {
		probFilterAnnotate := strings.Split(probeFiltersAnnotate, ",")
		for _, probeFilter := range probFilterAnnotate {
			probeFilters = append(probeFilters, "region: " + probeFilter)
		}
	}

	proto, err := utils.GetAnnotationValue(ing.Annotations, "protocol")
	if err != nil{
		proto = os.Getenv("PINGDOM_PROTOCOL")
	}
	customPath, _ := utils.GetAnnotationValue(ing.Annotations, "custom-path")

	var newPingdomCheck = newPingdomCheck{}
	var pingdomUrl = pClient.endpoint
	var method = "POST"
	var jsonValue []byte

	pc := pingdomCheck{
		Name: ing.Name,
		Host: ing.Spec.Rules[0].Host,
		CustomPath: customPath,
		Integrationids: integrationids,
		ProbeFilters: probeFilters,
	}

	if checkID := getCheckID(ing.Name); checkID != "" {
		pingdomUrl += "/" + checkID
		method = "PUT"
		jsonValue, _ = json.Marshal(pc)
	} else {
		newPingdomCheck.pingdomCheck = pc
		newPingdomCheck.Type = proto
		jsonValue, _ = json.Marshal(newPingdomCheck)
	}

	log.Println(string(jsonValue))
	sendPingdomRequest(method, pingdomUrl, jsonValue)
}

 func deleteCheck(ingName string){
	 if checkID := getCheckID(ingName); checkID != "" {
	 	message := `{"message": "Deletion of the check `+ingName+` was successful!"}`
		 sendPingdomRequest("DELETE", pClient.endpoint + "/" + checkID, []byte(message))
	 } else {
		 log.Printf("\nCannot find %s, delete failed\n", ingName)
	 }
 }

func getCheckID(checkName string) string {
	log.Printf("Sending request to pingdom to get check id")
	client := &http.Client{}
	req, _ := http.NewRequest("GET", pClient.endpoint, nil)
	req.Header.Set("Authorization", "Bearer " + pClient.token)
	res, err := client.Do(req)

	if err != nil {
		log.Printf("git check id resutl")
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	var pingdomChecks responsePingdomCheck
	if err := json.Unmarshal(body, &pingdomChecks); err != nil {
		panic(err)
	}
	for _, pc := range pingdomChecks.Checks{
		if pc.Name == checkName{
			return strconv.Itoa(pc.ID)
		}
	}
	return ""
}

 func sendPingdomRequest(method, url string, body []byte){
	 log.Printf("\nSending request to pingdom %s\n", url)

	 req := &http.Request{}
	 req, _ = http.NewRequest(method, url, bytes.NewBuffer(body))
	 req.Header.Set("Authorization", "Bearer " + pClient.token)
	 req.Header.Set("Content-Type", "application/json")

	 client := &http.Client{}
	 res, err := client.Do(req)

	 if err != nil {
		 log.Fatalln(err)
	 }
	 body, err = ioutil.ReadAll(res.Body)
	 if err != nil {
		 log.Fatalln(err)
	 }

	 if err := json.Unmarshal(body, &res); err != nil {
		 panic(err)
	 }

	 switch res.StatusCode {
	 case 200:
		 log.Printf("Succefully %s request on check %s", method, body)
		 break
	 case 404:
		 log.Printf("Not found!")
		 break
	 default:
		 log.Printf("Got response: %s", res.Status)
	 }
 }