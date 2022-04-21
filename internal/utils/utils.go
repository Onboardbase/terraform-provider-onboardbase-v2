package utils

import (
	"errors"
)

func Parseresult(resultData map[string]interface{}) (string, error) {
	generalProjects := resultData["generalPublicProjects"].(map[string]interface{})
	projectList := generalProjects["list"].([]interface{})
	if len(projectList) < 1 {
		return "", errors.New("Project not found")
	}
	var activeEnvironment map[string]interface{}

	if len(projectList) > 0 {
		projectData := projectList[0].(map[string]interface{})
		projectEnvironments := projectData["publicEnvironments"].(map[string]interface{})
		projectEnvironmentsData := projectEnvironments["list"].([]interface{})
		if len(projectEnvironmentsData) == 0 {
			return "", errors.New("Environment not found")
		}
		for _, environment := range projectEnvironmentsData {
			val, _ := environment.(map[string]interface{})

			if val["title"].(string) == "development" {
				activeEnvironment = val
				break
			}
		}
	}
	if activeEnvironment == nil {
		return "", errors.New("Environment not found")
	}
	projectEnvironmentsKeys := activeEnvironment["key"].(string)
	return projectEnvironmentsKeys, nil
}
