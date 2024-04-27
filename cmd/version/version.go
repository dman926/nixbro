package version

import (
	"github.com/dman926/nixbro/cmd"

	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Nixbro",
	Long:  `All software has versions. This is Nixbro's`,
	Run:   cmd.Version,
}
