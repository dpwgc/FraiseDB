package base

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func InitLog() error {
	err := CreatePath(Config().Store.Path)
	if err != nil {
		return err
	}
	logFile, err := os.OpenFile(fmt.Sprintf("%s/runtime-%s.log", Config().Store.Path, strings.Split(time.Now().String(), " ")[0]), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	LogHandler = log.Default()
	return nil
}
