package minishell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	gops "github.com/mitchellh/go-ps"
)

type Minishell struct{}

func New() *Minishell {
	return &Minishell{}
}

func (ms *Minishell) Execute(ctx context.Context, query string) {
	query = strings.TrimSpace(query)
	if query == "" {
		return
	}
	parts := strings.Split(query, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	if len(parts) == 1 {
		ms.executeSingle(ctx, parts[0])
		return
	}
	ms.executePipe(parts)
}

func (ms *Minishell) executeSingle(ctx context.Context, command string) {
	commandSlice := strings.Split(command, " ")
	if len(commandSlice) == 0 {
		return
	}
	for i := 1; i < len(commandSlice); i++ {
		value := commandSlice[i]
		if len(value) > 0 && []rune(value)[0] == '&' {
			commandSlice[i] = os.Getenv(string([]rune(value)[1:]))
		}
	}
	cmdName := commandSlice[0]
	switch cmdName {
	case "pwd":
		pwd()
	case "cd":
		if len(commandSlice) > 1 {
			cd(commandSlice[1])
		}
	case "echo":
		echo(commandSlice[1:])
	case "ps":
		ps()
	case "kill":
		if len(commandSlice) > 1 {
			pid, err := strconv.Atoi(commandSlice[1])
			if err == nil {
				kill(pid)
			} else {
				fmt.Println(err.Error())
			}
		}
	default:
		c := exec.CommandContext(ctx, commandSlice[0], commandSlice[1:]...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return
		}
	}
}

func (ms *Minishell) executePipe(parts []string) {
	numCmds := 0
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			numCmds++
		}
	}
	if numCmds < 2 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(numCmds)

	rpipe := make([]io.ReadCloser, numCmds)
	wpipe := make([]io.WriteCloser, numCmds-1)
	for i := 0; i < numCmds-1; i++ {
		var err error
		rpipe[i+1], wpipe[i], err = os.Pipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "pipe failed:", err)
			return
		}
	}

	cmdIdx := 0
	cmds := make([]*exec.Cmd, numCmds)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		commandSlice := strings.Split(part, " ")
		if len(commandSlice) == 0 {
			continue
		}
		cmd := exec.Command(commandSlice[0], commandSlice[1:]...)
		cmds[cmdIdx] = cmd
		if cmdIdx == 0 {
			cmd.Stdin = nil
		} else {
			cmd.Stdin = rpipe[cmdIdx]
		}
		if cmdIdx < numCmds-1 {
			cmd.Stdout = wpipe[cmdIdx]
		} else {
			cmd.Stdout = os.Stdout
		}
		cmd.Stderr = os.Stderr

		go func(idx int, c *exec.Cmd) {
			defer wg.Done()
			if err := c.Start(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			if idx < numCmds-1 {
				_ = wpipe[idx].Close()
			}
			_ = c.Wait()
		}(cmdIdx, cmd)
		cmdIdx++
	}

	wg.Wait()
	for _, w := range wpipe {
		_ = w.Close()
	}
	for i := 1; i < numCmds; i++ {
		_ = rpipe[i].Close()
	}
}

func cd(dir string) {
	if err := os.Chdir(dir); err != nil {
		fmt.Println(err.Error())
	}
}

func pwd() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(wd)
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
