package runner

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/avraam311/go-minishell/internal/minishell"
)

type App struct {
	minishell *minishell.Minishell
}

func New(ms *minishell.Minishell) *App {
	return &App{minishell: ms}
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

	scanner := bufio.NewScanner(os.Stdin)
	wd, _ := os.Getwd()

	for {
		fmt.Print(wd + "> ")
		if !scanner.Scan() {
			if err := scanner.Err(); err == io.EOF || err == nil {
				fmt.Println("\nexit")
				break
			}
		}

		query := scanner.Text()
		if query == "\\quit" {
			break
		}

		a.minishell.Execute(ctx, query)
		wd, _ = os.Getwd()
	}
}
