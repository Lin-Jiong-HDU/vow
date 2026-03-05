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
