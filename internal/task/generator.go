package task

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Lin-Jiong-HDU/vow/internal/config"
	"github.com/Lin-Jiong-HDU/vow/internal/done"
	"github.com/Lin-Jiong-HDU/vow/internal/word"
)

const (
	DefaultTaskDir = "tasks"
)

// DefaultTaskDirPath 返回默认任务目录路径
func DefaultTaskDirPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(home, ".vow", DefaultTaskDir), nil
}

// GenerateTasks 生成指定天数的单词栈
func GenerateTasks(wordList *word.WordList, doneList *done.DoneList, cfg *config.Config, days int) ([]*DailyTask, error) {
	// 移除已完成的单词
	wordList.RemoveWords(doneList.Words)

	// 根据配置决定是否打乱
	if cfg.Shuffle {
		wordList.Shuffle()
	}

	// 检查单词数量是否足够
	totalWords := wordList.Count()
	requiredWords := cfg.DailyWordCount * days
	if totalWords < requiredWords {
		return nil, fmt.Errorf("not enough words: have %d, need %d", totalWords, requiredWords)
	}

	// 生成任务
	tasks := make([]*DailyTask, days)
	wordIndex := 0

	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, i).Format("2006-01-02")
		task := NewDailyTask(date)

		// 分配单词
		for j := 0; j < cfg.DailyWordCount && wordIndex < totalWords; j++ {
			task.Words = append(task.Words, wordList.Words[wordIndex])
			wordIndex++
		}

		tasks[i] = task
	}

	return tasks, nil
}

// SaveTasks 保存所有任务到文件
func SaveTasks(tasks []*DailyTask, taskDir string) error {
	for _, task := range tasks {
		filename := task.Date + "-tasks.json"
		path := filepath.Join(taskDir, filename)

		if err := task.Save(path); err != nil {
			return fmt.Errorf("failed to save task for %s: %w", task.Date, err)
		}
	}
	return nil
}
