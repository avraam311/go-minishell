package minishell

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

type Command struct {
	Name string
	Args []string
}

type Minishell struct {
	Commands []*Command
}

func New() *Minishell {
	return &Minishell{
		Commands: []*Command{},
	}
}

func (m *Minishell) ParseLine(line string) {
	for _, part := range strings.Split(line, "|") {
		fields := strings.Fields(part)
		if len(fields) == 0 {
			continue
		}
		m.Commands = append(m.Commands, &Command{Name: fields[0], Args: fields[1:]})
	}
}

func (m *Minishell) IsBuiltin(cmd *Command) bool {
	builtins := map[string]bool{
		"cd":   true,
		"pwd":  true,
		"echo": true,
		"kill": true,
		"ps":   true,
	}
	return builtins[cmd.Name]
}

func (m *Minishell) ExecuteCmd(cmd *Command) int {
	if m.IsBuiltin(cmd) {
		switch cmd.Name {
		case "cd":
			err := m.cmdCd(cmd.Args)
			if err != nil {
				fmt.Println(err.Error())
				return 1
			}
		case "pwd":
			m.cmdPwd()
		case "echo":
			m.cmdEcho(cmd.Args)
		case "kill":
			err := m.cmdKill(cmd.Args)
			if err != nil {
				fmt.Println(err.Error())
				return 1
			}
		case "ps":
			m.cmdPs()
		default:
			fmt.Printf("Unknown builtin command '%s'\n", cmd.Name)
			return 1
		}
		return 0
	}

	command := exec.Command(cmd.Name, cmd.Args...)
	out, err := command.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", cmd.Name, err.Error())
		return 1
	}
	fmt.Print(string(out))
	return 0
}

func (m *Minishell) cmdCd(args []string) error {
	if len(args) != 1 {
		fmt.Println("Usage: cd <directory>")
		return nil
	}
	dir := args[0]
	err := os.Chdir(dir)
	if err != nil {
		fmt.Printf("Error changing directory to %s\n", dir)
	}
	return err
}

func (m *Minishell) cmdPwd() {
	cwd, _ := os.Getwd()
	fmt.Println(cwd)
}

func (m *Minishell) cmdEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func (m *Minishell) cmdKill(args []string) error {
	if len(args) != 1 {
		fmt.Println("Usage: kill <PID>")
		return nil
	}
	pidStr := args[0]
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		fmt.Printf("Invalid PID: %v\n", pidStr)
		return err
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("No such process with PID %d\n", pid)
		return err
	}
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		fmt.Printf("Failed to send signal to PID %d\n", pid)
	}
	return err
}

func (m *Minishell) cmdPs() {
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(string(output))
}
