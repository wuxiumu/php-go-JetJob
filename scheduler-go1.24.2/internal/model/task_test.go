package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTaskModelFields(t *testing.T) {
	now := time.Now()
	task := Task{
		Name:      "unit",
		Type:      "shell",
		Command:   "echo ok",
		Schedule:  "",
		Status:    "active",
		Retry:     1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	assert.Equal(t, "unit", task.Name)
	assert.Equal(t, "shell", string(task.Type))
	assert.Equal(t, 1, task.Retry)
}
