package os

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

/*
	(Parameterize) Clone config

	(Parameterize) Prompt for:
		Hostname
		Drive args (Show prompt of current disks and context for selection)

	(Parameterize toggle) Run disko to format disks (default)
		`sudo nix --experimental-features "nix-command flakes" run github:nix-community/disko -- --mode disko $DISKO PATH [--arg device '"/dev/nvme0nX"' --arg device2 '"/dev/sdX"']`
		!Update config to match!

	Generate hardware-configuration.nix
		`nixos-generate-config --no-filesystems --root /mnt --dir $FlakeLocation/hosts/$Hostname/`

	Attempt mount for install if disko not ran

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

type diskSelection struct {
	Arg string
	Val string
}

func RunInstall() tea.Msg {
	// PLACEHOLDERS
	remoteUrl := "https://github.com/dman926/nixos-config"
	configAbsPath, direrr := os.Getwd()
	hostname := "neutron"
	if direrr != nil {
		return InstallReturn{}
	}

	err := cloneConfig(remoteUrl, configAbsPath)
	if err != nil {
		return InstallReturn{}
	}
	disks := discPrompt(configAbsPath, hostname)
	err = runDisko(configAbsPath, hostname, disks)
	if err != nil {
		return InstallReturn{}
	}
	err = generateHardwareConfig(configAbsPath, hostname)
	if err != nil {
		return InstallReturn{}
	}
	err = scaffoldPersist(configAbsPath, hostname)
	if err != nil {
		return InstallReturn{}
	}
	err = installNixos(configAbsPath, hostname)
	if err != nil {
		return InstallReturn{}
	}
	err = scaffoldSecrets(configAbsPath, hostname)
	if err != nil {
		return InstallReturn{}
	}
	err = installNixos(configAbsPath, hostname)
	if err != nil {
		return InstallReturn{}
	}
	reboot()

	return InstallReturn{}
}

func cloneConfig(remoteUrl string, configAbsPath string) (err error) {
	cloneCmd := exec.Command("git", "clone", remoteUrl, configAbsPath)
	err = cloneCmd.Run()

	return
}

func discPrompt(configAbsPath string, hostname string) []diskSelection {
	// TODO
	return []diskSelection{}
}

func runDisko(configAbsPath string, hostname string, disks []diskSelection) (err error) {
	baseArgs := []string{
		"sudo", "nix", "--experimental-features", "nix-command flakes", "run", "github:nix-community/disko", "--", "--mode", "disko", configAbsPath + "/hosts/" + hostname + "/disks.nix",
	}
	combinedArgs := make([]string, len(disks)*3)
	for i := 0; i < len(disks); i++ {
		combinedArgs[i*3] = "--arg"
		combinedArgs[i*3+1] = disks[i].Arg
		combinedArgs[i*3+2] = "\"" + disks[i].Val + "\""
	}
	combinedArgs = append(baseArgs[1:], combinedArgs...)

	diskoCmd := exec.Command(baseArgs[0], combinedArgs...)
	err = diskoCmd.Run()

	return
}

func generateHardwareConfig(configAbsPath string, hostname string) (err error) {
	generateCmd := exec.Command("nixos-generate-config", "--no-filesystems", "--root", "/mnt", "--dir", fmt.Sprintf("%s/hosts/%s", configAbsPath, hostname))
	err = generateCmd.Run()

	return
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
