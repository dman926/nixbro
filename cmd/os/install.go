package os

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	osRootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install NixOS",
	Long:  `Install NixOS with my NixOS config`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Install")

		/*
			(Parameterize) Clone config

			(Parameterize) Prmpt for:
				Hostname
				Drive args (Show prompt of current disks and context for selection)

			(Parameterize toggle) Run disko to format disks (default) / Attempt mount for install
				`sudo nix --experimental-features "nix-command flakes" run github:nix-community/disko -- --mode disko $DISKO PATH [--arg device '"/dev/nvme0nX"' --arg device2 '"/dev/sdX"']`
				!Update config to match!

			Scaffold folders
				Scaffold /persist folders with correct ownership and permissions

			Install nixos
				`nixos-install --root /mnt --flake $FlakeLocation#$Hostname`

			Secrets
				Add machine's age public fingerprint to .sops.yaml and update secret files (nixos/secrets.yaml | nixos/users/*\/secrets.yaml)
					`nix-shell -p ssh-to-age --run 'cat /etc/ssh/ssh_host_ed25519_key.pub | ssh-to-age'`

			Install nixos again to get secrets working (namely user passwords)
				`nixos-install --root /mnt --flake $FlakeLocation#$Hostname`

			Reboot

			First user login prompt for final setup
				SSH / GPG keys
		*/
	},
}
