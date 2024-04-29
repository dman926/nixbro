package util

import (
	"bytes"
	"fmt"
	"os/exec"
)

// NixFlakeEval evaluates a Nix expression in the context of a Nix flake.
// Useful for retreiving config values
func NixFlakeEval(flakeRootAbsPath string, expr string) (output bytes.Buffer, err error) {

	// TODO: allow for array to be passed for multiple commands and outputs
	// want to avoid spawning the repl a ton
	// Just involves needing to split the output buffer into a slice of buffers to match each expr position

	cmd := exec.Command("nix", "repl", "--quiet")

	cmd.Dir = flakeRootAbsPath

	var inBuffer bytes.Buffer
	fmt.Fprintln(&inBuffer, ":lf .#")
	fmt.Fprintln(&inBuffer, ":p "+expr)
	fmt.Fprintln(&inBuffer, ":q")
	cmd.Stdin = &inBuffer

	cmd.Stdout = &output

	// Run the command
	err = cmd.Run()

	return
}
