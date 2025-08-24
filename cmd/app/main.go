package app

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/avraam311/go-minishell/internal/minishell"
)

type App struct {
	minishell *minishell.Minishell
}

func New(ms *minishell.Minishell) *App {
	return &App{
		minishell: ms,
	}
}

func (a *App) Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = line[:len(line)-1]
		a.minishell.ParseLine(line)
		status := 0
		for i, cmd := range a.minishell.Commands {
			if i+1 < len(a.minishell.Commands) {
				r, w, err := os.Pipe()
				if err != nil {
					fmt.Println("Pipe failed:", err)
					os.Exit(1)
				}
				defer r.Close()
				defer w.Close()

				child := exec.Command(cmd.Name, cmd.Args...)
				child.Stdout = w
				go func(child *exec.Cmd) {
					err := child.Run()
					if err != nil {
						fmt.Println("cannot run the command")
						os.Exit(1)
					}
				}(child)

				w.Close()

				cmd = &minishell.Command{Name: "cat", Args: []string{"-"}}
				cmdIn, _ := os.CreateTemp("", "")
				err = os.WriteFile(cmdIn.Name(), make([]byte, 0), 0600)
				if err != nil {
					fmt.Println("cannot write into temporary file")
					os.Exit(1)
				}
				cmdIn.Close()
				cmd.Args = append(cmd.Args, cmdIn.Name())
			}

			status += a.minishell.ExecuteCmd(cmd)
		}

		if status != 0 {
			fmt.Println("Last command exited with non-zero exit code.")
		}
	}
}
