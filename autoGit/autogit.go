package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var dirs = []string{
	"/home/db/DingBen/MyNote",
	"/home/db/DingBen/Myscript",
}

// 判断是否有文件被更改过
func gitstatus(path string) bool {

	// 查看git状态
	cmd := exec.Command("git", "status")
	cmd.Dir = path

	gitstatus, err := cmd.Output()
	if err != nil {
		fmt.Println("git status err:", err)
	}

	return strings.Contains(string(gitstatus), "Changes not staged for commit")
}

func NowTime() string {
	// 返回当前时间，用于git commit标识
	return time.Now().Format("2006-01-02 15:04:05")
}

func autoGit(path string) {
	if gitstatus(path) {
		addCmd := exec.Command("git", "add", ".")
		commitCmd := exec.Command("git", "commit", "-m", NowTime())
		pushCmd := exec.Command("git", "push", "origin", "main")

		addCmd.Dir = path
		commitCmd.Dir = path
		pushCmd.Dir = path

		if err := addCmd.Run(); err != nil {
			fmt.Println("git add err:", err)
			return
		}

		if err := commitCmd.Run(); err != nil {
			fmt.Println("git commit err:", err)
			return
		}

		if err := pushCmd.Run(); err != nil {
			fmt.Println("push origin master:", err)
			return
		}
	}

}

func main() {
	for _, dir := range dirs {
		autoGit(dir)
	}
}
