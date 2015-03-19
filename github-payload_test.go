package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushEventIsMatch(t *testing.T) {
	c := Criteria{
		Event:      "push",
		Owner:      "Pica9",
		Repository: "cd-core",
	}
	e := PushEvent{}
	e.Repository.Name = "cd-core"
	e.Repository.Owner.Name = "Pica9"
	assert.Equal(t, e.IsMatch(c), true, "Expected a match")
	e.Repository.Owner.Name = "someone else"
	assert.Equal(t, e.IsMatch(c), false, "Expected a match")
	e.Repository.Name = "otherrepo"
	e.Repository.Owner.Name = "Pica9"
	assert.Equal(t, e.IsMatch(c), false, "Expected a match")
}
