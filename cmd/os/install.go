package os

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/dman926/nixbro/util"

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

type InstallReturn struct{}

type diskSelection struct {
	Arg         string
	Val         string
	Description string
	Current     string
}

type InitInstallReturn []string

func InitInstall(remoteUrl string, configAbsPath string) tea.Cmd {
	return func() tea.Msg {
		// TODO: check if config is there first. FF if it is
		err := cloneConfig(remoteUrl, configAbsPath)
		if err != nil {
			return InstallReturn{}
		}

		files, err := os.ReadDir(configAbsPath + "/hosts")
		if err != nil {
			log.Fatal(err)
		}

		// Filter the results to only include directories
		hosts := []string{}
		for _, file := range files {
			if file.IsDir() {
				hosts = append(hosts, file.Name())
			}
		}

		return InitInstallReturn(hosts)
	}
}

type DisksReturn []diskSelection

func CollectDisks(configAbsPath string, hostname string) tea.Cmd {
	return func() tea.Msg {
		disks, err := collectDisks(configAbsPath, hostname)
		if err != nil {
			return InstallReturn{}
		}

		return DisksReturn(disks)
	}
}

func RunInstall(configAbsPath string, hostname string, disks []diskSelection) tea.Cmd {
	return func() tea.Msg {
		err := runDisko(configAbsPath, hostname, disks)
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
}

func cloneConfig(remoteUrl string, configAbsPath string) (err error) {
	cloneCmd := exec.Command("git", "clone", remoteUrl, configAbsPath)
	err = cloneCmd.Run()

	return
}

func collectDisks(configAbsPath string, hostname string) (disks []diskSelection, err error) {
	disksFileAbsPath := configAbsPath + "/hosts/" + hostname + "/disks.nix"
	disksFile, err := os.Open(disksFileAbsPath)
	if err != nil {
		return
	}

	stat, err := disksFile.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := stat.Size()

	// Read the file content
	disksFileContent := make([]byte, fileSize)
	_, err = disksFile.Read(disksFileContent)
	if err != nil {
		log.Fatal(err)
	}

	// Find bounds of params
	startParams := 0
	endParams := 0
	for cursor := 0; cursor < len(disksFileContent); cursor++ {
		if disksFileContent[cursor] == '{' {
			startParams = cursor
		}
		if disksFileContent[cursor] == '}' {
			endParams = cursor
			break
		}
	}

	params := strings.Split(string(disksFileContent[startParams+1:endParams]), ",")
	for _, param := range params {
		parts := strings.SplitN(param, "?", 1)
		if !strings.Contains(strings.ToLower(parts[0]), "device") {
			// Not a device parameter
			continue
		}

		disk := diskSelection{
			Arg: strings.TrimSpace(parts[0]),
		}
		if len(parts) == 2 {
			parts[1] = strings.TrimSpace(parts[1])
			targetStartWord := "throw"
			if len(parts[1]) > len(targetStartWord) && parts[1][0:len(targetStartWord)] == targetStartWord {
				startQuote := strings.Index(parts[1], "\"")
				endQuote := strings.LastIndex(parts[1], "\"")
				disk.Description = strings.TrimSpace(parts[1][startQuote+1 : endQuote])
			}
		}
		disks = append(disks, disk)
	}

	// TODO: fetch currently set values from configuration.nix in same dir

	return
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
	generateCmd := exec.Command("nixos-generate-config", "--no-filesystems", "--root", "/mnt", "--dir", configAbsPath+"/hosts/"+hostname)
	err = generateCmd.Run()

	return
}

func scaffoldPersist(configAbsPath string, hostname string) (err error) {
	elemAt := func(wrapping string, index int) string { return "builtins.elemAt (" + wrapping + ") " + string(index) }
	attrValues := func(wrapping string) string { return "builtins.attrValues (" + wrapping + ")" }

	systemPersistenceExpr := "outputs.nixosConfigurations." + hostname + ".config.environment.persistence"
	systemPersistenceCmd := elemAt(attrValues("builtins.mapAttrs (name: value: (map (val: name + val.directory) value.directories)) "+systemPersistenceExpr), 0)
	systemPersistence, err := util.NixFlakeEval(configAbsPath, systemPersistenceCmd)
	if err != nil {
		return
	}

	users := []string{} // TODO: get list of users from /hosts/<hostname>/home/*.nix names

	homePersistences := []bytes.Buffer{}
	for _, user := range users {
		homePersistenceExpr := "outputs.homeConfigurations." + user + "@" + hostname + ".config.home.persistence"
		homePersistencesCmd := elemAt(attrValues("builtins.mapAttrs (name: value: (map (val: name + \"/\" + (if val ? \"directory\" then val.directory else val)) value.directories)) "+homePersistenceExpr), 0)
		homePersistence, err := util.NixFlakeEval(configAbsPath, homePersistencesCmd)
		if err != nil {
			// TODO
			continue
		}

		homePersistences = append(homePersistences, homePersistence)
	}

	// TODO scaffold directories listed in each persistence
	// Both persistences are set up as array of absolute paths (["/persist/system/...", ...], ["/persist/home/USER/...", ...])
	// Ensure ownership and permissions for home persistences are set properly

	return
}

func installNixos(configAbsPath string, hostname string) (err error) {
	installCmd := exec.Command("nixos-install", "--root", "/mnt", "--flake", configAbsPath, "#", hostname)
	err = installCmd.Run()

	return
}

func scaffoldSecrets(configAbsPath string, hostname string) (err error) {
	ageCmd := exec.Command("nix-shell", "-p", "ssh-to-age", "--run", "cat /etc/ssh/ssh_host_ed25519_key.pub | ssh-to-age")
	output, err := ageCmd.Output()
	if err != nil {
		return
	}
	agePublicKey := strings.TrimRight(string(output), "\n")

	// TODO: Add public key to .sops.yaml and update secret files (nixos/secrets.yaml | nixos/users/*\/secrets.yaml)

	return
}

func reboot() {
	// TODO
}
