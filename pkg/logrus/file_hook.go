package logrus

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type fileHook struct {
	LogPath string
}

func NewFileHook(path string) *fileHook {
	return &fileHook{
		LogPath: path,
	}
}

func (hook *fileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *fileHook) Fire(entry *logrus.Entry) error {
	fileName := fmt.Sprintf("%v/log_%v.json", hook.LogPath, time.Now().Format("2006-01-02"))
	fmt.Println(fileName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	logData, _ := entry.Bytes()

	file.WriteString(string(logData))

	return nil
}
