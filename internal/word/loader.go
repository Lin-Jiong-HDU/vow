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
