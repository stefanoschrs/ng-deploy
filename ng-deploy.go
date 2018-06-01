package main

import (
	"os"
	"log"
	"fmt"
	"flag"
	"bufio"
	"os/exec"
	"strings"
	"encoding/json"
	"time"
)

type Configuration struct {
	Ssh struct {
		User string `json:"user"`
		Path string `json:"path"`
		Host string `json:"host"`
		Key string `json:"key"`
	} `json:"ssh"`
	Backup struct {
		Prefix string `json:"prefix"`
		Path string `json:"path"`
	} `json:"backup"`
}

type CliFlags struct {
	Build bool
	Sync bool
	Backup bool
	Env string
}

var config Configuration
var cliFlags CliFlags

/**
 * Commands
 */
func buildCmd() {
	if _, err := os.Stat("package.json"); os.IsNotExist(err) {
		log.Fatal("'package.json' could not be found")
	}

	cmd := exec.Command("npm", "run", "build:" + cliFlags.Env)

	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanBytes)

	for scanner.Scan() {
		m := scanner.Text()
		fmt.Print(m)
	}
	err := cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func backupCmd() {
	lastSlash := strings.LastIndex(config.Ssh.Path, "/")
	basePath := config.Ssh.Path[0:lastSlash]
	appPath := config.Ssh.Path[lastSlash + 1:len(config.Ssh.Path)]

	cmd := exec.Command("ssh",
		"-i", config.Ssh.Key,
		config.Ssh.User + "@" + config.Ssh.Host,
		fmt.Sprintf("tar -zcvf %s/%s.%d.tgz -C %s %s",
			config.Backup.Path,
			config.Backup.Prefix,
			time.Now().Unix(),
			basePath,
			appPath,
		),
	)

	stderr, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
}

func syncCmd() {
	if _, err := os.Stat("dist"); os.IsNotExist(err) {
		log.Fatal("'dist' could not be found")
	}

	cmd := exec.Command("rsync",
		"-arvP",
		"-e", "ssh -i " + config.Ssh.Key,
		"dist/",
		config.Ssh.User + "@" + config.Ssh.Host + ":" + config.Ssh.Path,
	)

	stderr, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
}

/**
 * Misc
 */
func parseFlags()  {
	cliFlags = CliFlags{}

	// TODO: Use FlagSet to have both short, long options
	flag.StringVar(&cliFlags.Env, "env", "", "Environment")
	flag.BoolVar(&cliFlags.Build, "build", false, "Build project")
	flag.BoolVar(&cliFlags.Sync, "sync", false, "Sync files")
	flag.BoolVar(&cliFlags.Backup, "backup", false, "Backup target")

	flag.Parse()

	if cliFlags.Env == "" {
		log.Fatal("Missing '--env'")
	}
}

func loadConfiguration() {
	file, err := os.Open(".ng-deploy.json")
	defer file.Close()

	if err != nil && os.IsExist(err) {
		log.Fatal(err)
	}

	if file == nil {
		return
	}

	var f map[string]Configuration
	decoder := json.NewDecoder(file)
	decoder.Decode(&f)

	config = f[cliFlags.Env]

	// make sure there is no trailing /
	config.Ssh.Path = strings.TrimSuffix(config.Ssh.Path, "/")
	config.Backup.Path = strings.TrimSuffix(config.Backup.Path, "/")
}

/**
 * Main
 */
func main() {
	parseFlags()
	loadConfiguration()

	if cliFlags.Build {
		buildCmd()
	}

	if cliFlags.Backup {
		backupCmd()
	}

	if cliFlags.Sync {
		syncCmd()
	}
}
