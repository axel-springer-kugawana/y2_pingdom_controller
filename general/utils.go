package general

import (
	"log"
	"strconv"
	"strings"

)

func GetAnnotationValue (annotations map[string]string, contains string) string{
	for k, v := range annotations {
		if strings.Contains(k, contains){
			log.Printf("Found anntation %s with value %s", k, v)
			return v
		}
	}
	log.Printf("Annotation with %s not found", contains)
	return ""
}

func StringToInt(str string) int{
	var result int
	var err error
	if result, err = strconv.Atoi(str); err != nil {
		log.Printf("Cannot convert String to Int %s" , str)
	}
	return result
}

func StringToBool(str string) bool{
	var result bool
	var err error
	if result, err = strconv.ParseBool(str); err != nil {
		log.Printf("Cannot convert String to Boolen %s" , str)
		return false
	}
	return result
}

func GetBoolPointer(b bool) *bool{
	return &b
}
