package runner

import (
	"bufio"
	"fmt"
	"os"

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
	wd, _ := os.Getwd()
	fmt.Print(wd + "> ")
	for scanner := bufio.NewScanner(os.Stdin); scanner.Scan(); fmt.Print(wd + "> ") {
		if query := scanner.Text(); query != "\\quit" {
			a.minishell.Execute(query)
		} else {
			break
		}
		wd, _ = os.Getwd()
	}
}
