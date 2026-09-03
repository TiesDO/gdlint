package main

import (
	"github.com/TiesDO/gdlint/cmd"
	"github.com/TiesDO/gdlint/core"
	"github.com/TiesDO/gdlint/rules"
)

func init() {
	rules.RegisterAll(&core.DefaultRuleRegistry)
}

func main() {
	cmd.Execute()
}
