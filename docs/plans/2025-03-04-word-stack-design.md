# 单词栈模块设计文档

**日期**: 2025-03-04
**状态**: 已批准

## 概述

实现单词表加载、已完成单词管理和单词栈生成功能。加载单词表后自动生成未来10天的单词栈。

## 模块结构

```
internal/
├── config/
│   └── config.go       # 添加 Shuffle 字段
├── word/
│   ├── word.go         # Word 结构体 + 辅助方法
│   └── loader.go       # 加载单词表
├── done/
│   ├── done.go         # Done 相关结构体
│   └── loader.go       # 加载已完成单词
└── task/
    ├── task.go         # Task 结构体
    └── generator.go    # 生成单词栈
```

## 1. config 模块修改

### 修改内容

在 `Config` 结构体中添加 `Shuffle` 字段：

```go
type Config struct {
    DailyWordCount int  `json:"dailyWordCount"`
    Shuffle        bool `json:"shuffle"`  // 是否打乱单词顺序，默认 false
}
```

## 2. word 模块

### 数据结构

```go
// Word 表示一个单词
type Word struct {
    Word    string
    Meaning string
    Example string
}

// WordList 表示单词表
type WordList struct {
    Name  string
    Words []Word
}
```

### API

```go
// 构造函数
func NewWordList() *WordList

// 辅助方法
func (wl *WordList) Count() int
func (wl *WordList) Shuffle()                    // Fisher-Yates 洗牌
func (wl *WordList) RemoveWords(words []string)  // 移除指定单词

// 加载
const DefaultWordListFile = "vow.json"
func DefaultWordListPath() (string, error)
func LoadWordList(path string) (*WordList, error)
```

### 行为

- 文件不存在时返回空 WordList
- JSON 解析失败返回错误
- 路径为空时使用 `~/.vow/vow.json`

## 3. done 模块

### 数据结构

```go
// DoneList 已完成单词总列表
type DoneList struct {
    UpdateTime time.Time
    Words      []string
}

// DailyDone 每日已完成单词
type DailyDone struct {
    Date  string
    Words []string
}
```

### API

```go
// 构造函数
func NewDoneList() *DoneList

// 辅助方法
func (dl *DoneList) Contains(word string) bool
func (dl *DoneList) AddWords(words []string)
func (dl *DoneList) GetWordSet() map[string]bool

// 加载
const DefaultDoneFile = "done.json"
const DefaultDailyDoneDir = "done"
func DefaultDonePath() (string, error)
func LoadDoneList(path string) (*DoneList, error)
```

### 行为

- 文件不存在时返回空 DoneList
- 路径为空时使用 `~/.vow/done.json`

## 4. task 模块

### 数据结构

```go
// DailyTask 每日单词栈
type DailyTask struct {
    Date  string
    Words []Word
}
```

### API

```go
// 构造函数
func NewDailyTask(date string) *DailyTask

// 辅助方法
func (dt *DailyTask) Count() int
func (dt *DailyTask) Save(path string) error

// 生成
const DefaultTaskDir = "tasks"
func DefaultTaskDirPath() (string, error)
func GenerateTasks(wordList *WordList, doneList *DoneList, cfg *Config, days int) ([]*DailyTask, error)
```

### 生成逻辑

1. 从 wordList 中移除 doneList 中的单词
2. 如果 cfg.Shuffle 为 true，打乱剩余单词
3. 按顺序将单词分配给指定天数
4. 每天分配 cfg.DailyWordCount 个单词
5. 返回生成的 DailyTask 列表

### 文件格式

每日任务保存为 `~/.vow/tasks/{YYYY-MM-DD}-tasks.json`：

```json
{
  "date": "2025-03-05",
  "words": [
    {
      "word": "abandon",
      "meaning": "v. 放弃",
      "example": "Example sentence"
    }
  ]
}
```

## 数据流

```
LoadConfig()
    ↓
LoadWordList() → 加载单词表
    ↓
LoadDoneList() → 加载已完成单词
    ↓
wordList.RemoveWords(doneList.Words) → 移除已完成单词
    ↓
if cfg.Shuffle → wordList.Shuffle()
    ↓
GenerateTasks(wordList, doneList, cfg, 10) → 生成10天单词栈
    ↓
dailyTask.Save() → 保存到 ~/.vow/tasks/{YYYY-MM-DD}-tasks.json
```

## 文件结构

```
~/.vow/
├── config.json           # 配置文件
├── vow.json              # 单词表
├── done.json             # 已完成单词总列表
├── done/                 # 每日已完成记录
│   └── {YYYY-MM-DD}-done.json
└── tasks/                # 每日单词栈
    └── {YYYY-MM-DD}-tasks.json
```

## 常量定义

```go
const (
    DefaultWordListDir  = ".vow"
    DefaultWordListFile = "vow.json"
    DefaultDoneFile     = "done.json"
    DefaultDailyDoneDir = "done"
    DefaultTaskDir      = "tasks"
)
```

## 错误处理

- 文件不存在：返回空列表（word, done）
- JSON 解析失败：返回错误
- 目录创建失败：返回错误
- 文件写入失败：返回错误
