package logger

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type (
	File struct {
		*logrus.Logger
		file *os.File
	}
)

func NewFile(filePath string) (*File, error) {

	logDir := filepath.Dir(filePath)
	_, err := os.Stat(logDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(file)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     false,
	})

	return &File{
		Logger: logger,
		file:   file,
	}, nil
}

func (l *File) Close() error {
	return l.file.Close()
}
