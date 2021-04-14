package controller

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/russellcardullo/go-pingdom/pingdom"
	extensions "k8s.io/api/extensions/v1beta1"
	utils "pingdom_controller/general"
)

// PingdomEngine struct
type PingdomEngine struct {
	addIngress    chan *extensions.Ingress
	updateIngress chan *extensions.Ingress
	deleteIngress chan string
}

// NewPingdomEngine is running a channel for each event
func NewPingdomEngine() *PingdomEngine {
	return &PingdomEngine{
		addIngress:    make(chan *extensions.Ingress),
		updateIngress: make(chan *extensions.Ingress),
		deleteIngress: make(chan string),
	}
}

var client, _ = pingdom.NewClientWithConfig(pingdom.ClientConfig{
	APIToken: os.Getenv("PINGDOM_TOKEN"),
})

// Run will start PingdomEngine's channels
func (p *PingdomEngine) Run() {
	for {
		select {
		case ing := <-p.addIngress:
			createNewCheck(ing)
		case ing := <-p.updateIngress:
			updateCheck(ing)
		case ingName := <-p.deleteIngress:
			deleteCheck(ingName)
		}
	}
}

func createPingdomCheckType(ing *extensions.Ingress) pingdom.HttpCheck {
	check := pingdom.HttpCheck{
		Name:              ing.Name,
		Hostname:          ing.Spec.Rules[0].Host,
		Resolution:        utils.StringToInt(utils.GetAnnotationValue(ing.Annotations, "resolution")),
		Paused:            utils.StringToBool(utils.GetAnnotationValue(ing.Annotations, "paused")),
		Encryption:        utils.StringToBool(utils.GetAnnotationValue(ing.Annotations, "encryption")),
		Url:               utils.GetAnnotationValue(ing.Annotations, "custom-path"),
		Port:              utils.StringToInt(utils.GetAnnotationValue(ing.Annotations, "port")),
		IntegrationIds:    extractAndBuildArrayOfIntegers(ing, "integrationids"),
		Tags:              utils.GetAnnotationValue(ing.Annotations, "port"),
		ProbeFilters:      utils.GetAnnotationValue(ing.Annotations, "probe-filters"),
		TeamIds:           extractAndBuildArrayOfIntegers(ing, "teamids"),
		VerifyCertificate: utils.GetBoolPointer(utils.StringToBool(utils.GetAnnotationValue(ing.Annotations, "verify-certificate"))),
	}
	return check
}

func createNewCheck(ing *extensions.Ingress) {
	if checkID := getCheckID(ing.Name); checkID != "" {
		log.Printf("\nCheck with the name %s is already exist. Starting update operation\n", ing.Name)
		updateCheck(ing)
		return
	}

	pct := createPingdomCheckType(ing)
	// Create a new http check
	check, err := client.Checks.Create(&pct)
	if err != nil {
		log.Println("Error", err)
	}
	log.Println("Created check:", check) // {ID, Name}
}

func extractAndBuildArrayOfIntegers(ing *extensions.Ingress, annotate string) []int {
	userInput := utils.GetAnnotationValue(ing.Annotations, annotate)
	var userInputAsArray = []int{}
	if userInput == "" {
		log.Printf("No %s input found in annotations list", annotate)
	} else {
		userInputSplit := strings.Split(userInput, ",")
		for _, input := range userInputSplit {
			res, err := strconv.Atoi(input)
			if err != nil {
				log.Fatal("Cannot convert to int in Integrationids")
				break
			}
			userInputAsArray = append(userInputAsArray, res)
		}
	}
	return userInputAsArray
}

func updateCheck(ing *extensions.Ingress) {
	pct := createPingdomCheckType(ing)
	if checkID := getCheckID(pct.Name); checkID != "" {
		msg, err := client.Checks.Update(utils.StringToInt(checkID), &pct)
		if err != nil {
			log.Printf("\nError while trying to update %s check", ing.Name)
		} else {
			log.Printf(msg.Message)
		}
	} else {
		log.Printf("\nCannot find %s check, update failed. Starting create operation\n", ing.Name)
		createNewCheck(ing)
	}
}

func deleteCheck(ingName string) {
	if checkID := getCheckID(ingName); checkID != "" {
		msg, err := client.Checks.Delete(utils.StringToInt(checkID))
		if err != nil {
			log.Printf("\nError while trying to delete %s check", ingName)
		} else {
			log.Printf(msg.Message)
		}
	} else {
		log.Printf("\nCannot find %s check, delete failed\n", ingName)
	}
}

func getCheckID(checkName string) string {
	checks, _ := client.Checks.List()
	for _, check := range checks {
		if check.Name == checkName {
			return strconv.Itoa(check.ID)
		}
	}
	log.Printf("\nCannot find %s check\n", checkName)
	return ""
}
