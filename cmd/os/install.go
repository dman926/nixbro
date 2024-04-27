package os

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	osRootCmd.AddCommand(installCmd)

	// installCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	// installCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
	// installCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
	// installCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
	installCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	// viper.BindPFlag("author", installCmd.PersistentFlags().Lookup("author"))
	// viper.BindPFlag("projectbase", installCmd.PersistentFlags().Lookup("projectbase"))
	viper.BindPFlag("useViper", installCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache")
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install NixOS",
	Long:  `Install NixOS with my NixOS config`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Install")
	},
}
