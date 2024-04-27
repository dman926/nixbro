package os

import (
	"fmt"

	"github.com/dman926/nixbro/cmd"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

func init() {
	cmd.RootCmd.AddCommand(osRootCmd)
}

var osRootCmd = &cobra.Command{
	Use:   "os",
	Short: "os stuff...",
	Long:  `...`,
	Run:   Help,
}

func Help(cmd *cobra.Command, args []string) {
	fmt.Println("`nix os ...` help")
}
