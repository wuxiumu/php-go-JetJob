package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("JETJOB_TEST_CONFIG", "../../config/config.yaml")
	LoadConfig(os.Getenv("JETJOB_TEST_CONFIG"))
	assert.NotEmpty(t, Cfg.MySQL.DSN)
}
