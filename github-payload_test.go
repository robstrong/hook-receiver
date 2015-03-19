package main

import (
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
	assert.Equal(t, e.IsMatch(c), false, "Expected a match")

	c.PushParams.Branch = "mast*"
	assert.Equal(t, e.IsMatch(c), true, "Expected a match")

	c.Event = "release"
	assert.Equal(t, e.IsMatch(c), false, "Expected a match")
}
