package runner

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan os.Signal, 1)
	go func() {
		signal.Notify(ch, os.Interrupt)
		<-ch
		cancel()
	}()
	wd, _ := os.Getwd()
	fmt.Print(wd + "> ")
	for scanner := bufio.NewScanner(os.Stdin); scanner.Scan(); fmt.Print(wd + "> ") {
		if query := scanner.Text(); query != "\\quit" {
			a.minishell.Execute(ctx, query)
		} else {
			break
		}
		wd, _ = os.Getwd()
	}
}
