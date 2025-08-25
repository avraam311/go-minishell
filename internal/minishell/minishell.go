package minishell

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	gops "github.com/mitchellh/go-ps"
)

type Minishell struct {
}

func New() *Minishell {
	return &Minishell{}
}

func (ms *Minishell) Execute(query string) {
	commands := strings.Split(query, " | ")
	for _, command := range commands {
		commandSlice := strings.Split(command, " ")
		switch commandSlice[0] {
		case "pwd":
			pwd()
		case "cd":
			cd(commandSlice[1])
		case "echo":
			echo(commandSlice[1:])
		case "ps":
			ps()
		case "kill":
			pid, err := strconv.Atoi(commandSlice[1])
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			kill(pid)
		}
	}

}

func cd(dir string) {
	if err := os.Chdir(dir); err != nil {
		fmt.Println(err.Error())
	}
}

func pwd() string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	fmt.Println(wd)
	return wd
}

func echo(args []string) {
	for _, word := range args {
		fmt.Print(word + " ")
	}
	fmt.Println()
}

func ps() {
	if pcs, err := gops.Processes(); err != nil {
		fmt.Println(err.Error())
	} else {
		for _, pc := range pcs {
			fmt.Println(pc.Pid(), pc.Executable())
		}
	}
}

func kill(pid int) {
	p, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := p.Kill(); err != nil {
		fmt.Println(err.Error())
		return
	}
}
