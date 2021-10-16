package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"time"
)

// TODO:使用json文件保存路径
var dirs = []string{
	"/home/db/DingBen/MyNote",
}

// 返回修改过的文件名和未跟踪的文件名
func gitstatus(path string) (modified []string, untracked []string) {

	// 查看git状态
	cmd := exec.Command("git", "status")
	cmd.Dir = path

	gitstatus, err := cmd.Output()
	if err != nil {
		fmt.Println("git status err:", err)
	}

	modified = findModified(string(gitstatus))
	untracked = findUntracked(string(gitstatus))
	return
}

func findModified(status string) (modified []string) {
	reg := regexp.MustCompile(`modified:\s{3}(.*)\n`)
	if reg == nil {
		fmt.Println("regex err")
		return
	}
	result := reg.FindAllStringSubmatch(status, -1)
	for _, text := range result {
		modified = append(modified, text[1])
	}
	return
}

func findUntracked(status string) (untracked []string) {
	reg1 := regexp.MustCompile(`use "git add <file>..." to include in what will be committed.\n((.*\n){1,})no changes`)
	reg2 := regexp.MustCompile(`\t{0,}(.*)`)

	var tmp string
	if len(reg1.FindAllStringSubmatch(status, -1)) > 0 {
		tmp = reg1.FindAllStringSubmatch(status, -1)[0][1]
	} else {
		return
	}
	for _, v := range reg2.FindAllStringSubmatch(tmp, -1) {
		untracked = append(untracked, v[1])
	}
	return
}

func NowTime() string {
	// 返回当前时间，用于git commit标识
	return time.Now().Format("2006-01-02 15:04:05")
}

func addAndCommit(filepath string, dir string, comment string) {
	addCmd := exec.Command("git", "add", filepath)
	commitCmd := exec.Command("git", "commit", "-m", comment+NowTime())
	addCmd.Dir = dir
	commitCmd.Dir = dir

	if err := addCmd.Run(); err != nil {
		fmt.Println("git add", filepath, "err", err)
		return
	}

	if err := commitCmd.Run(); err != nil {
		fmt.Println("git commit", filepath, "err", err)
		return
	}
}

func autoGit(dir string) {
	modified, untracked := gitstatus(dir)

	// git add修改过的文件，并commit 修改于什么时间
	if modified != nil {
		for _, v := range modified {
			addAndCommit(v, dir, "修改于")
		}
	}

	// git add未跟踪的文件，并commit 创建于什么时间
	if untracked != nil {
		for _, v := range untracked {
			addAndCommit(v, dir, "创建于")
		}
	}

	// 自动提交到远程仓库
	pushCmd := exec.Command("git", "push", "origin", "main")
	pushCmd.Dir = dir

	if err := pushCmd.Run(); err != nil {
		fmt.Println("push origin master:", err)
		return
	}

}

func main() {
	for _, dir := range dirs {
		autoGit(dir)
	}
}
