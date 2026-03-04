package word

import "math/rand"

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
