# 单词栈模块实施计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 实现单词表加载、已完成单词管理和未来10天单词栈自动生成功能

**Architecture:** 分四个模块实现 - config 添加 Shuffle 字段、word 模块加载单词表、done 模块管理已完成单词、task 模块生成单词栈。数据流：加载配置 → 加载单词表 → 加载已完成单词 → 去重 → 可选打乱 → 生成10天任务 → 保存文件。

**Tech Stack:** Go 1.x, encoding/json, time 包, Fisher-Yates 洗牌算法

---

### Task 1: 修改 config.go 添加 Shuffle 字段

**Files:**
- Modify: `internal/config/config.go`

**Step 1: 修改 Config 结构体**

在 `internal/config/config.go` 的 `Config` 结构体中添加 `Shuffle` 字段：

```go
type Config struct {
    DailyWordCount int  `json:"dailyWordCount"`
    Shuffle        bool `json:"shuffle"`
}
```

**Step 2: 更新默认值常量**

在 `New()` 函数中设置默认值：

```go
func New() *Config {
    return &Config{
        DailyWordCount: DefaultDailyWordCount,
        Shuffle:        false,  // 默认不打乱
    }
}
```

**Step 3: 验证编译通过**

Run: `go build ./...`
Expected: 无错误

**Step 4: 测试默认配置**

Run: `go run cmd/main.go`
Expected: 应该显示 `Shuffle: false`

**Step 5: Commit**

```bash
git add internal/config/config.go
git commit -m "feat(config): 添加 Shuffle 字段支持单词乱序"
```

---

### Task 2: 创建 word 模块 - word.go

**Files:**
- Create: `internal/word/word.go`

**Step 1: 创建文件并定义结构体**

```bash
mkdir -p internal/word
```

创建 `internal/word/word.go`:

```go
package word

// Word 表示一个单词及其释义和例句
type Word struct {
    Word    string `json:"word"`
    Meaning string `json:"meaning"`
    Example string `json:"example"`
}

// WordList 表示单词表
type WordList struct {
    Name  string  `json:"name"`
    Words []Word  `json:"words"`
}

// NewWordList 创建一个空的单词表
func NewWordList() *WordList {
    return &WordList{
        Name:  "",
        Words: []Word{},
    }
}

// Count 返回单词数量
func (wl *WordList) Count() int {
    return len(wl.Words)
}

// Shuffle 使用 Fisher-Yates 算法打乱单词顺序
func (wl *WordList) Shuffle() {
    n := len(wl.Words)
    for i := n - 1; i > 0; i-- {
        j := rand.Intn(i + 1)
        wl.Words[i], wl.Words[j] = wl.Words[j], wl.Words[i]
    }
}

// RemoveWords 从单词表中移除指定的单词
func (wl *WordList) RemoveWords(words []string) {
    wordSet := make(map[string]bool)
    for _, w := range words {
        wordSet[w] = true
    }

    var newWords []Word
    for _, w := range wl.Words {
        if !wordSet[w.Word] {
            newWords = append(newWords, w)
        }
    }
    wl.Words = newWords
}

// GetWordsAsStrings 返回单词列表的字符串形式
func (wl *WordList) GetWordsAsStrings() []string {
    words := make([]string, len(wl.Words))
    for i, w := range wl.Words {
        words[i] = w.Word
    }
    return words
}
```

**Step 2: 添加 rand 包导入**

在文件顶部添加：

```go
import "math/rand"
```

**Step 3: 验证编译通过**

Run: `go build ./internal/word/...`
Expected: 无错误

**Step 4: Commit**

```bash
git add internal/word/word.go
git commit -m "feat(word): 添加 Word 和 WordList 结构体及辅助方法"
```

---

### Task 3: 创建 word 模块 - loader.go

**Files:**
- Create: `internal/word/loader.go`

**Step 1: 创建 loader.go 文件**

```go
package word

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

const (
    DefaultWordListDir  = ".vow"
    DefaultWordListFile = "vow.json"
)

// DefaultWordListPath 返回默认单词表路径
func DefaultWordListPath() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("failed to get user home directory: %w", err)
    }
    return filepath.Join(home, DefaultWordListDir, DefaultWordListFile), nil
}

// LoadWordList 加载单词表
func LoadWordList(path string) (*WordList, error) {
    if path == "" {
        var err error
        path, err = DefaultWordListPath()
        if err != nil {
            return nil, err
        }
    }

    fileData, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            // 文件不存在，返回空列表
            return NewWordList(), nil
        }
        return nil, fmt.Errorf("failed to read word list file: %w", err)
    }

    var wordList WordList
    if err := json.Unmarshal(fileData, &wordList); err != nil {
        return nil, fmt.Errorf("failed to parse word list: %w", err)
    }

    return &wordList, nil
}
```

**Step 2: 验证编译通过**

Run: `go build ./internal/word/...`
Expected: 无错误

**Step 3: Commit**

```bash
git add internal/word/loader.go
git commit -m "feat(word): 添加 LoadWordList 函数支持加载单词表"
```

---

### Task 4: 创建 done 模块 - done.go

**Files:**
- Create: `internal/done/done.go`

**Step 1: 创建文件并定义结构体**

```bash
mkdir -p internal/done
```

创建 `internal/done/done.go`:

```go
package done

import "time"

// DoneList 已完成单词总列表
type DoneList struct {
    UpdateTime time.Time `json:"updateTime"`
    Words      []string  `json:"words"`
}

// DailyDone 每日已完成单词
type DailyDone struct {
    Date  string   `json:"date"`
    Words []string `json:"words"`
}

// NewDoneList 创建一个空的已完成单词列表
func NewDoneList() *DoneList {
    return &DoneList{
        UpdateTime: time.Now(),
        Words:      []string{},
    }
}

// Contains 检查单词是否已完成
func (dl *DoneList) Contains(word string) bool {
    for _, w := range dl.Words {
        if w == word {
            return true
        }
    }
    return false
}

// AddWords 添加已完成单词
func (dl *DoneList) AddWords(words []string) {
    dl.Words = append(dl.Words, words...)
    dl.UpdateTime = time.Now()
}

// GetWordSet 返回单词集合
func (dl *DoneList) GetWordSet() map[string]bool {
    wordSet := make(map[string]bool)
    for _, w := range dl.Words {
        wordSet[w] = true
    }
    return wordSet
}

// Count 返回已完成单词数量
func (dl *DoneList) Count() int {
    return len(dl.Words)
}
```

**Step 2: 验证编译通过**

Run: `go build ./internal/done/...`
Expected: 无错误

**Step 3: Commit**

```bash
git add internal/done/done.go
git commit -m "feat(done): 添加 DoneList 和 DailyDone 结构体及辅助方法"
```

---

### Task 5: 创建 done 模块 - loader.go

**Files:**
- Create: `internal/done/loader.go`

**Step 1: 创建 loader.go 文件**

```go
package done

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

const (
    DefaultDoneFile     = "done.json"
    DefaultDailyDoneDir = "done"
)

// DefaultDonePath 返回默认已完成单词列表路径
func DefaultDonePath() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("failed to get user home directory: %w", err)
    }
    return filepath.Join(home, ".vow", DefaultDoneFile), nil
}

// LoadDoneList 加载已完成单词列表
func LoadDoneList(path string) (*DoneList, error) {
    if path == "" {
        var err error
        path, err = DefaultDonePath()
        if err != nil {
            return nil, err
        }
    }

    fileData, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            // 文件不存在，返回空列表
            return NewDoneList(), nil
        }
        return nil, fmt.Errorf("failed to read done list file: %w", err)
    }

    var doneList DoneList
    if err := json.Unmarshal(fileData, &doneList); err != nil {
        return nil, fmt.Errorf("failed to parse done list: %w", err)
    }

    return &doneList, nil
}
```

**Step 2: 验证编译通过**

Run: `go build ./internal/done/...`
Expected: 无错误

**Step 3: Commit**

```bash
git add internal/done/loader.go
git commit -m "feat(done): 添加 LoadDoneList 函数支持加载已完成单词"
```

---

### Task 6: 创建 task 模块 - task.go

**Files:**
- Create: `internal/task/task.go`

**Step 1: 创建文件并定义结构体**

```bash
mkdir -p internal/task
```

创建 `internal/task/task.go`:

```go
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
    Date  string         `json:"date"`
    Words []word.Word    `json:"words"`
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
```

**Step 2: 验证编译通过**

Run: `go build ./internal/task/...`
Expected: 无错误

**Step 3: Commit**

```bash
git add internal/task/task.go
git commit -m "feat(task): 添加 DailyTask 结构体及 Save 方法"
```

---

### Task 7: 创建 task 模块 - generator.go

**Files:**
- Create: `internal/task/generator.go`

**Step 1: 创建 generator.go 文件**

```go
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
```

**Step 2: 验证编译通过**

Run: `go build ./internal/task/...`
Expected: 无错误

**Step 3: Commit**

```bash
git add internal/task/generator.go
git commit -m "feat(task): 添加 GenerateTasks 和 SaveTasks 函数"
```

---

### Task 8: 更新 main.go 集成所有模块

**Files:**
- Modify: `cmd/main.go`

**Step 1: 重写 main.go**

```go
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
```

**Step 2: 验证编译通过**

Run: `go build`
Expected: 无错误

**Step 3: Commit**

```bash
git add cmd/main.go
git commit -m "feat(main): 集成所有模块，实现单词栈生成功能"
```

---

### Task 9: 创建测试用单词表

**Files:**
- Create: 手动创建测试文件

**Step 1: 创建测试单词表**

创建 `~/.vow/vow.json` (手动创建或使用命令):

```json
{
  "name": "test-words",
  "words": [
    {
      "word": "abandon",
      "meaning": "v. 放弃；抛弃",
      "example": "He abandoned his family."
    },
    {
      "word": "ability",
      "meaning": "n. 能力；才能",
      "example": "She has the ability to solve problems."
    },
    {
      "word": "able",
      "meaning": "adj. 能够的；有能力的",
      "example": "I am able to swim."
    },
    {
      "word": "about",
      "meaning": "prep. 关于；大约",
      "example": "Tell me about yourself."
    },
    {
      "word": "above",
      "meaning": "prep. 在...之上",
      "example": "The plane flew above the clouds."
    },
    {
      "word": "abroad",
      "meaning": "adv. 在国外",
      "example": "She went abroad last year."
    },
    {
      "word": "absence",
      "meaning": "n. 缺席；不在",
      "example": "His absence was noticed."
    },
    {
      "word": "absent",
      "meaning": "adj. 缺席的",
      "example": "He was absent from school."
    },
    {
      "word": "absolute",
      "meaning": "adj. 绝对的；完全的",
      "example": "I have absolute confidence."
    },
    {
      "word": "absorb",
      "meaning": "v. 吸收；吸引",
      "example": "Plants absorb sunlight."
    },
    {
      "word": "abstract",
      "meaning": "adj. 抽象的",
      "example": "Truth is an abstract concept."
    },
    {
      "word": "abundant",
      "meaning": "adj. 丰富的；充裕的",
      "example": "The forest has abundant wildlife."
    },
    {
      "word": "abuse",
      "meaning": "v./n. 滥用；虐待",
      "example": "Don't abuse your power."
    },
    {
      "word": "academic",
      "meaning": "adj. 学术的",
      "example": "She pursued an academic career."
    },
    {
      "word": "academy",
      "meaning": "n. 学院；学会",
      "example": "He graduated from the military academy."
    },
    {
      "word": "accelerate",
      "meaning": "v. 加速；促进",
      "example": "The car began to accelerate."
    },
    {
      "word": "accept",
      "meaning": "v. 接受；同意",
      "example": "I accept your apology."
    },
    {
      "word": "access",
      "meaning": "n. 通路；访问权",
      "example": "Do you have internet access?"
    },
    {
      "word": "accident",
      "meaning": "n. 事故；意外",
      "example": "He was injured in an accident."
    },
    {
      "word": "accomplish",
      "meaning": "v. 完成；实现",
      "example": "She accomplished her goal."
    }
  ]
}
```

**Step 4: 运行程序测试**

Run: `go run cmd/main.go`
Expected: 输出显示加载了 20 个单词，生成了 10 个任务

**Step 5: 验证生成的文件**

Run: `ls -la ~/.vow/tasks/`
Expected: 应该有 10 个 JSON 文件，格式为 `YYYY-MM-DD-tasks.json`

Run: `cat ~/.vow/tasks/$(date +%Y-%m-%d)-tasks.json`
Expected: 应该显示今天的任务，包含 2 个单词（默认 dailyWordCount=20，但测试只有 20 个单词，所以每天 2 个）

---

### Task 10: 添加 rand seed 初始化

**Files:**
- Modify: `internal/word/word.go`

**Step 1: 修改 Shuffle 方法**

在 `word.go` 顶部的 import 中添加:

```go
import (
    "math/rand"
    "time"
)
```

修改 `Shuffle` 方法:

```go
// Shuffle 使用 Fisher-Yates 算法打乱单词顺序
func (wl *WordList) Shuffle() {
    rand.Seed(time.Now().UnixNano())  // 初始化随机种子
    n := len(wl.Words)
    for i := n - 1; i > 0; i-- {
        j := rand.Intn(i + 1)
        wl.Words[i], wl.Words[j] = wl.Words[j], wl.Words[i]
    }
}
```

**Step 2: 验证编译通过**

Run: `go build ./...`
Expected: 无错误

**Step 3: Commit**

```bash
git add internal/word/word.go
git commit -m "fix(word): 添加随机种子初始化"
```

---

### Task 11: 测试乱序功能

**Step 1: 修改配置测试乱序**

编辑 `~/.vow/config.json`:

```json
{
  "dailyWordCount": 2,
  "shuffle": true
}
```

**Step 2: 删除旧任务重新生成**

Run: `rm -rf ~/.vow/tasks/ && go run cmd/main.go`
Expected: 生成 10 个任务，单词顺序是打乱的

**Step 3: 验证乱序效果**

Run: `cat ~/.vow/tasks/$(date +%Y-%m-%d)-tasks.json | head -20`
Expected: 每次运行单词顺序应该不同

**Step 4: 最终验证**

Run: `go build && go test ./...`
Expected: 所有测试通过

---

## 验收标准

1. ✅ `config` 模块支持 `Shuffle` 字段
2. ✅ `word` 模块可以加载单词表，文件不存在时返回空列表
3. ✅ `done` 模块可以加载已完成单词列表
4. ✅ `task` 模块可以生成未来 10 天的单词栈
5. ✅ 单词栈会自动排除已完成的单词
6. ✅ 当 `shuffle=true` 时，单词顺序被打乱
7. ✅ 所有模块编译通过，程序正常运行
