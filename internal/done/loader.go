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
