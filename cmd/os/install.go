package os

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

/*
	(Parameterize) Clone config

	(Parameterize) Prompt for:
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

	First user login prompt for final setup (this will come later)
		SSH / GPG keys
*/

type InstallReturn struct {
}

func RunInstall() tea.Msg {
	return InstallReturn{}
}

func cloneConfig(remoteUrl string, configAbsPath string) (err error) {
	cloneCmd := exec.Command("git", "clone", remoteUrl, configAbsPath)
	err = cloneCmd.Run()

	return
}

func discPrompt(configAbsPath string, hostname string) {
	// TODO
}

func scaffoldPersist(configAbsPath string, hostname string) (err error) {
	return
}

func installNixos(configAbsPath string, hostname string) (err error) {
	installCmd := exec.Command("nixos-install", "--root /mnt", "--flake \""+fmt.Sprintf("%s#%s", configAbsPath, hostname)+"\"")
	err = installCmd.Run()

	return
}

func scaffoldSecrets(configAbsPath string, hostname string) (err error) {
	// TODO
	return
}

func reboot() {
	// TODO
}
