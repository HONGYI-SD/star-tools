package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	log.Println("watch dog start!")
	go doCheck()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	switch <-quit {
	case syscall.SIGINT:
		log.Println("interrupt")
	case syscall.SIGTERM:
		log.Println("terminated")
	default:
		log.Println("default")
	}
}

func doCheck() {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			detectCmd := "ps aux|grep -w '../dong'|grep -v 'grep'|awk '{print $2}'"
			cmd, err := Exec(detectCmd)
			if err != nil {
				log.Println("watch dog err, continue")
				continue
			}
			pid := strings.ReplaceAll(cmd, "\n", "")
			if pid != "" {
				log.Println("app is running")
				continue
			}

			//重新拉起进程
			runCmd := "nohup ../dong > tmp.log 2>&1 &"
			_, error := Exec(runCmd)
			if error == nil {
				log.Println("App started!")
			} else {
				log.Println(error)
			}
		}
	}
}

func Exec(args string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", args)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", errors.New(stderr.String())
	}
	return strings.Trim(stdout.String(), "\n"), nil
}
