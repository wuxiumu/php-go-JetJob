package app

import (
	"encoding/json"
	"errors"
	"io"
	"jetjob/internal/model"
	"jetjob/internal/storage"
	"net/http"
	"strings"
	"time"
)

// CreateTask 创建任务
func CreateTask(t *model.Task) error {
	// 可添加业务校验、默认值赋值等
	if t.Name == "" || t.Type == "" {
		return errors.New("任务名称和类型不能为空")
	}
	if t.Retry == 0 {
		t.Retry = 3
	}
	return storage.DB.Create(t).Error
}

// ListTasks 获取所有任务
func ListTasks() ([]model.Task, error) {
	var ts []model.Task
	err := storage.DB.Find(&ts).Error
	return ts, err
}

// ExecuteTask 根据任务类型分发执行
func ExecuteTask(t *model.Task) (string, error) {
	switch t.Type {
	case model.TASK_TYPE_SHELL:
		return executeShell(t.Command, t.Params)
	case model.TASK_TYPE_HTTP:
		return executeHttp(t.Command, t.Params)
	case model.TASK_TYPE_FILE:
		return executeFileTask(t.Command, t.Params)
	default:
		return "", errors.New("未知任务类型")
	}
}

// executeHttp 执行HTTP类型任务
func executeHttp(url, params string) (string, error) {
	var p struct {
		Method  string            `json:"method"`
		Headers map[string]string `json:"headers"`
		Body    string            `json:"body"`
	}
	if err := json.Unmarshal([]byte(params), &p); err != nil {
		return "", errors.New("HTTP任务参数解析失败: " + err.Error())
	}
	if p.Method == "" {
		p.Method = "GET"
	}
	req, err := http.NewRequest(p.Method, url, strings.NewReader(p.Body))
	if err != nil {
		return "", errors.New("创建HTTP请求失败: " + err.Error())
	}
	for k, v := range p.Headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("HTTP请求失败: " + err.Error())
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("读取HTTP响应失败: " + err.Error())
	}
	return string(respBytes), nil
}

// executeShell 执行Shell类型任务（请根据实际安全性和功能扩展）
func executeShell(command, params string) (string, error) {
	// TODO: 安全实现 shell 执行，可以用 os/exec，参数校验
	return "模拟执行shell：" + command, nil
}

// executeFileTask 执行文件相关任务（上传/下载/处理等）
func executeFileTask(command, params string) (string, error) {
	// TODO: 实现你的文件任务逻辑
	return "模拟执行文件任务：" + command, nil
}
