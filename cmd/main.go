package main

import (
	"fmt"
	"log"

	"github.com/Lin-Jiong-HDU/vow/internal/config"
	"github.com/Lin-Jiong-HDU/vow/internal/done"
	"github.com/Lin-Jiong-HDU/vow/internal/task"
	"github.com/Lin-Jiong-HDU/vow/internal/word"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("Config loaded: DailyWordCount=%d, Shuffle=%v\n", cfg.DailyWordCount, cfg.Shuffle)

	// 加载单词表
	wordList, err := word.LoadWordList("")
	if err != nil {
		log.Fatalf("Failed to load word list: %v", err)
	}
	fmt.Printf("Word list loaded: %d words\n", wordList.Count())

	// 加载已完成单词
	doneList, err := done.LoadDoneList("")
	if err != nil {
		log.Fatalf("Failed to load done list: %v", err)
	}
	fmt.Printf("Done list loaded: %d words completed\n", doneList.Count())

	// 生成10天单词栈
	tasks, err := task.GenerateTasks(wordList, doneList, cfg, 10)
	if err != nil {
		log.Fatalf("Failed to generate tasks: %v", err)
	}

	// 保存任务
	taskDir, err := task.DefaultTaskDirPath()
	if err != nil {
		log.Fatalf("Failed to get task dir: %v", err)
	}

	if err := task.SaveTasks(tasks, taskDir); err != nil {
		log.Fatalf("Failed to save tasks: %v", err)
	}

	fmt.Printf("Generated %d tasks successfully!\n", len(tasks))
	for _, t := range tasks {
		fmt.Printf("  - %s: %d words\n", t.Date, t.Count())
	}
}
