package general

import (
	"log"
	"strconv"
	"strings"
)

// GetAnnotationValue extract the value from annotations that contain a certain string
func GetAnnotationValue (annotations map[string]string, contains string) string{
	for k, v := range annotations {
		if strings.Contains(k, contains){
			log.Printf("Found anntation %s with value %s", k, v)
			return v
		}
	}
	return ""
}

// StringToInt convert string to int
func StringToInt(str string) int{
	var result int
	var err error
	if result, err = strconv.Atoi(str); err != nil {
		log.Printf("Cannot convert String to Int %s" , str)
	}
	return result
}

// StringToBool convert a string to a boolean
func StringToBool(str string) bool{
	var result bool
	var err error
	if result, err = strconv.ParseBool(str); err != nil {
		log.Printf("Cannot convert String to Boolean %s" , str)
		return false
	}
	return result
}

// GetBoolPointer function will return a pointer for a boolean variable
func GetBoolPointer(b bool) *bool{
	return &b
}
