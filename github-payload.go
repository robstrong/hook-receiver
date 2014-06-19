package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Payload struct {
	Type       string
	Owner      string
	Repository string
}

func formatPayload(req *http.Request) (payload Payload, err error) {
	var jsonBody struct {
		Repository struct {
			Owner struct {
				Login string //these two hold the same data, but the push event format is different
				Name  string //from all the other events
			}
			Name string
		}
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		return
	}

	payload.Type = req.Header.Get("X-Github-Event")
	payload.Owner = jsonBody.Repository.Owner.Name
	if payload.Owner == "" {
		payload.Owner = jsonBody.Repository.Owner.Login
	}
	payload.Repository = jsonBody.Repository.Name
	json, err := json.Marshal(payload)
	fmt.Println("Payload: ", string(json))

	return
}
