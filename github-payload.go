package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

//go:generate gojson -input=github-events/push.json -name=PushEvent -o=push_event.go
//go:generate gojson -input=github-events/release.json -name=ReleaseEvent -o=release_event.go
//go:generate gojson -input=github-events/issues.json -name=IssuesEvent -o=issues_event.go

type Payload interface {
	IsMatch(Criteria) bool
	Type() string
}

func parsePayload(req *http.Request) (Payload, error) {
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return PushEvent{}, err
	}

	event := req.Header.Get("X-Github-Event")

	switch event {
	case "push":
		pushEvent := PushEvent{}
		err = json.Unmarshal(body, &pushEvent)
		return pushEvent, err
	case "release":
		releaseEvent := ReleaseEvent{}
		err = json.Unmarshal(body, &releaseEvent)
		return releaseEvent, err
	case "issues":
		issuesEvent := IssuesEvent{}
		err = json.Unmarshal(body, &issuesEvent)
		return issuesEvent, err
	default:
		return PushEvent{}, errors.New("invalid event type: " + event)
	}
}

func (e PushEvent) Type() string {
	return "push"
}
func (e PushEvent) Branch() string {
	if strings.HasPrefix(e.Ref, "refs/heads/") {
		return e.Ref[11:]
	}
	return ""
}
func (e PushEvent) Tag() string {
	if strings.HasPrefix(e.Ref, "refs/tags/") {
		return e.Ref[10:]
	}
	return ""
}
func (e PushEvent) IsMatch(c Criteria) bool {
	//check that types match
	if e.Type() != c.Event {
		return false
	}
	//check that owners match
	if c.Owner != "" && e.Repository.Owner.Name != c.Owner {
		return false
	}
	//check that repo names match
	if c.Repository != "" && e.Repository.Name != c.Repository {
		return false
	}
	//check that branch matches
	if c.PushParams.Branch != "" {
		matched, err := regexp.MatchString("refs/heads/"+c.PushParams.Branch, e.Ref)
		if err != nil {
			return false
		}
		if !matched {
			return false
		}
	}
	return true
}

func (e ReleaseEvent) Type() string {
	return "release"
}
func (e ReleaseEvent) IsMatch(c Criteria) bool {
	//check that types match
	if e.Type() != c.Event {
		return false
	}
	//check that owners match
	if c.Owner != "" && e.Repository.Owner.Login != c.Owner {
		return false
	}
	//check that repo names match
	if c.Repository != "" && e.Repository.Name != c.Repository {
		return false
	}
	//check that prerelease matches
	if c.ReleaseParams.Prerelease != nil && *c.ReleaseParams.Prerelease != e.Release.Prerelease {
		return false
	}
	return true
}

func (e IssuesEvent) Type() string {
	return "issues"
}
func (e IssuesEvent) IsMatch(c Criteria) bool {
	//check that types match
	if e.Type() != c.Event {
		return false
	}
	//check that owners match
	if c.Owner != "" && e.Repository.Owner.Login != c.Owner {
		return false
	}
	//check that repo names match
	if c.Repository != "" && e.Repository.Name != c.Repository {
		return false
	}
	return true
}
