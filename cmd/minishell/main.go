package main

import (
	"github.com/avraam311/go-minishell/cmd/runner"
	"github.com/avraam311/go-minishell/internal/minishell"
)

func main() {
	minishell := minishell.New()
	runner := runner.New(minishell)
	runner.Run()
}
