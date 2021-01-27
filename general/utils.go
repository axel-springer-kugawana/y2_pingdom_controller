package general

import (
	"errors"
	"log"
	"strings"
)

func GetAnnotationValue (annotations map[string]string, contains string) (string, error){
	for k, v := range annotations {
		if strings.Contains(k, contains){
			log.Printf("Found anntation %s with value %s", k, v)
			return v, nil
		}
	}
	return "", errors.New("Annotation with "+contains+" not found")
}
