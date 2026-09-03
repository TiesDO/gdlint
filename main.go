package main

import (
	"github.com/TiesDO/gdlint/cmd"

	_ "github.com/TiesDO/gdlint/rules" // This is a hack to ensure the init() functions of the rules run
)

func main() {
	cmd.Execute()
}
