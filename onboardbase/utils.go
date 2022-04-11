package onboardbase

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

func Parseresult(resultData map[string]interface{}) ([]string, error) {
	generalProjects := resultData["generalPublicProjects"].(map[string]interface{})
	projectList := generalProjects["list"].([]interface{})
	if len(projectList) < 1 {
		return nil, errors.New("Project not found")
	}
	var activeEnvironment map[string]interface{}

	if len(projectList) > 0 {
		projectData := projectList[0].(map[string]interface{})
		projectEnvironments := projectData["publicEnvironments"].(map[string]interface{})
		projectEnvironmentsData := projectEnvironments["list"].([]interface{})
		if len(projectEnvironmentsData) == 0 {
			return nil, errors.New("Environment not found")
		}
		for _, environment := range projectEnvironmentsData {
			val, _ := environment.(map[string]interface{})

			if val["title"].(string) == "development" {
				activeEnvironment = val
				break
			}
		}
	}
	projectEnvironmentsKeys := activeEnvironment["key"].(string)
	var secrets []string
	err := json.Unmarshal([]byte(projectEnvironmentsKeys), &secrets)
	if err != nil {
		return nil, err
	}
	return secrets, err
}

func DecryptSecrets(secrets []string, passcode string, secret string, client *http.Client) (string, error) {
	reqBody := map[string]interface{}{}
	reqBody["secrets"] = secrets
	reqBody["secret"] = secret
	reqBody["passcode"] = passcode
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://w8n9u7hcd9.execute-api.us-east-1.amazonaws.com/decoder", strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	r, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if r.StatusCode == 400 {
		return "", errors.New("Environment variable not found")
	}
	defer r.Body.Close()
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	err = json.Unmarshal(bodyBytes, &result)

	if err != nil {
		return "", err
	}
	if result["error"] != nil {
		return "", errors.New(result["error"].(string))
	}

	return result["value"].(string), err
}
