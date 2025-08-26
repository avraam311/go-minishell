package minishell

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	gops "github.com/mitchellh/go-ps"
)

type Minishell struct{}

func New() *Minishell {
	return &Minishell{}
}

func (ms *Minishell) Execute(ctx context.Context, query string) {
	stop := false
	commands := strings.Split(query, " | ")
	for _, command := range commands {
		if stop {
			break
		}
		commandSlice := strings.Split(command, " ")
		for i, value := range commandSlice[1:] {
			if []rune(value)[0] == '&' {
				commandSlice[i+1] = os.Getenv(string([]rune(value)[1:]))
			}
		}
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
		default:
			cmd := exec.Command(commandSlice[0], commandSlice[1:]...)
			out, err := cmd.Output()

			go func() {
				<-ctx.Done()
				stop = true
				err := cmd.Process.Signal(syscall.SIGINT)
				if err != nil {
					fmt.Println("error interrupting command:", err)
				}
			}()

			if err != nil {
				fmt.Println()
				continue
			}
			fmt.Print(string(out))
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
	fmt.Println("killed")
}
