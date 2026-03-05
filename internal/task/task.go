package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Lin-Jiong-HDU/vow/internal/word"
)

// DailyTask 每日单词栈
type DailyTask struct {
	Date  string      `json:"date"`
	Words []word.Word `json:"words"`
}

// NewDailyTask 创建一个空的每日任务
func NewDailyTask(date string) *DailyTask {
	return &DailyTask{
		Date:  date,
		Words: []word.Word{},
	}
}

// Count 返回单词数量
func (dt *DailyTask) Count() int {
	return len(dt.Words)
}

// Save 保存每日任务到文件
func (dt *DailyTask) Save(path string) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create task directory: %w", err)
	}

	// 序列化为 JSON
	data, err := json.MarshalIndent(dt, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write task file: %w", err)
	}

	return nil
}
