package os

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	osRootCmd.AddCommand(rebuildCmd)
}

var rebuildCmd = &cobra.Command{
	Use:   "rebuild",
	Short: "Rebuild NixOS",
	Long:  `Rebuild NixOS with the available NixOS config and any modifiers`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rebuild")
	},
}
