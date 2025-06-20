package app

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"jetjob/internal/model"
	"jetjob/internal/storage"
	"log"
	"os"
	"testing"
)

func setup() {
	// 指定 .env 文件的路径，假设 .env 文件在项目的根目录下
	err := godotenv.Load("../../.env") // 根据实际情况调整路径
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}

	dsn := os.Getenv("JETJOB_TEST_DSN")
	if dsn == "" {
		log.Fatal("Error getting JETJOB_TEST_DSN from .env file")
	}

	storage.InitDB(dsn) // 用测试库DSN

	storage.DB.AutoMigrate(&model.Task{})
}

func TestCreateListTask(t *testing.T) {
	setup()
	tk := model.Task{Name: "appcase", Type: "shell"}
	err := CreateTask(&tk)
	assert.NoError(t, err)

	list, err := ListTasks()
	assert.NoError(t, err)
	assert.True(t, len(list) > 0)
}
