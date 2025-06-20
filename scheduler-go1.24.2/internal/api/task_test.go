package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"jetjob/internal/model"
	"jetjob/internal/storage"
	"log"
	"net/http"
	"net/http/httptest"
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

func TestCreateTaskHandler(t *testing.T) {
	setup()
	router := gin.Default()
	router.POST("/api/tasks", CreateTaskHandler)

	task := model.Task{Name: "api_case", Type: "shell"}
	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}
