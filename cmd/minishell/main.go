package main

import (
	"github.com/avraam311/go-minishell/cmd/app"
	"github.com/avraam311/go-minishell/internal/minishell"
)

func main() {
	minishell := minishell.New()
	app := app.New(minishell)
	app.Run()
}
