package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushEventIsMatch(t *testing.T) {
	e := PushEvent{}
	e.Repository.Name = "cd-core"
	e.Repository.Owner.Name = "Pica9"
	e.Ref = "refs/heads/master"

	c := Criteria{
		Event: "push",
	}
	assert.Equal(t, e.IsMatch(c), true, "Expected a match")

	c.Owner = "Pica9"
	assert.Equal(t, e.IsMatch(c), true, "Expected a match")

	c.Repository = "cd-core"
	assert.Equal(t, e.IsMatch(c), true, "Expected a match")

	c.PushParams.Branch = "master"
	assert.Equal(t, e.IsMatch(c), true, "Expected a match")

	c.PushParams.Branch = "testing"
	assert.Equal(t, e.IsMatch(c), false, "Did not expect a match")

	c.PushParams.Branch = "mast*"
	assert.Equal(t, e.IsMatch(c), true, "Expected a match")

	c.Event = "release"
	assert.Equal(t, e.IsMatch(c), false, "Did not expect a match")

}

func TestPushEventIsMatchReleaseBranch(t *testing.T) {
	e := PushEvent{}
	e.Repository.Name = "cd-core"
	e.Repository.Owner.Name = "Pica9"
	e.Ref = "refs/heads/release-4.16.0"

	c := Criteria{
		Event:      "push",
		Owner:      "Pica9",
		Repository: "cd-core",
	}
	c.PushParams.Branch = "release-4.17.0"

	assert.Equal(t, e.IsMatch(c), false, "Did not expect a match")

	e.Ref = "refs/heads/release-4.17.0"
	assert.Equal(t, e.IsMatch(c), true, "Expected a match")
}

func TestParseConfigAndMatch(t *testing.T) {
	jsonReader := bytes.NewBufferString(`
		{
			"port": 8000,
			"rules": [
				{
					"command": "echo test",
					"criteria": [
						{
							"event": "push",
							"owner": "Pica9",
							"repository": "cd-core",
							"push_params": {
								"branch": "release-4.17.0"
							}
						}
					]
				}
			]
        }`)
	config := getConfigFromReader(jsonReader)

	assert.Equal(t, 1, len(config.Rules), "Rules not loaded properly")
	assert.Equal(t, "echo test", config.Rules[0].Command, "Command not parsed correctly")
	assert.Equal(t, 1, len(config.Rules[0].Criteria), "Criteria not parsed correctly")
	assert.Equal(t, "release-4.17.0", config.Rules[0].Criteria[0].PushParams.Branch, "Branch not parsed correctly")

	e := PushEvent{}
	e.Repository.Name = "cd-core"
	e.Repository.Owner.Name = "Pica9"
	e.Ref = "refs/heads/release-4.16.0"
	assert.Equal(t, e.IsMatch(config.Rules[0].Criteria[0]), false, "Did not expect a match")
	e.Ref = "refs/heads/release-4.17.0"
	assert.Equal(t, e.IsMatch(config.Rules[0].Criteria[0]), true, "Expected a match")
}
