package main

import (
	"github.com/dman926/nixbro/cmd"
	_ "github.com/dman926/nixbro/cmd/os"
	_ "github.com/dman926/nixbro/cmd/version"
)

func main() {
	cmd.Execute()
}
