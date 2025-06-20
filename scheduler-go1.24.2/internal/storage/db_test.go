// internal/storage/db_test.go
package storage

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"jetjob/internal/model"
	"os"
	"testing"
)

func TestInitDBAndCRUD(t *testing.T) {
	// 指定 .env 文件的路径，假设 .env 文件在项目的根目录下
	err := godotenv.Load("../../.env") // 根据实际情况调整路径
	if err != nil {
		t.Fatal("Error loading .env file: " + err.Error())
	}

	dsn := os.Getenv("JETJOB_TEST_DSN")
	if dsn == "" {
		t.Fatal("JETJOB_TEST_DSN environment variable is not set")
	}

	InitDB(dsn) // 用测试库DSN
	defer DB.Exec("DROP TABLE IF EXISTS tasks")

	// 自动迁移
	DB.AutoMigrate(&model.Task{}, &model.Node{})

	// Create
	task := model.Task{Name: "test", Type: model.TASK_TYPE_SHELL, Command: "echo test", Params: "{}", Schedule: "@every 1m", Trigger: "manual", Status: "active", Retry: 3}
	err = DB.Create(&task).Error
	assert.NoError(t, err)

	// Read
	var out model.Task
	err = DB.First(&out, "name = ?", "test").Error
	assert.NoError(t, err)
	assert.Equal(t, model.TASK_TYPE_SHELL, out.Type)

	// Update
	DB.Model(&out).Update("Status", "done")
	var updated model.Task
	DB.First(&updated, "id = ?", out.ID)
	assert.Equal(t, "done", updated.Status)

	// Delete
	DB.Delete(&model.Task{}, out.ID)
	var deleted model.Task
	r := DB.First(&deleted, "id = ?", out.ID)
	assert.Error(t, r.Error)
}
