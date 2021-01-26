package general

import (
	"errors"
	"strings"
)

func GetAnnotationValue (annotations map[string]string, contains string) (string, error){
	for k, v := range annotations {
		if strings.Contains(k, contains){
			return v, nil
		}
	}
	return "", errors.New("Annotation not found")
}
