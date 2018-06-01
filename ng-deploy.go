package main

import (
	"flag"
	"os"
	"log"
	"os/exec"
	"bufio"
	"fmt"
	"encoding/json"
	"strings"
)

type Configuration struct {
	Ssh struct {
		User string `json:"user"`
		Path string `json:"path"`
		Host string `json:"host"`
		Key string `json:"key"`
	} `json:"ssh"`
}

type CliFlags struct {
	Build bool
	Sync bool
	Backup bool
	Env string
}

var cliFlags CliFlags
var config Configuration
var projectDir string

func build() {
	if _, err := os.Stat("package.json"); os.IsNotExist(err) {
		log.Fatal("'package.json' could not be found")
	}

	cmd := exec.Command("npm", "run", "build:prod")

	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanBytes)

	for scanner.Scan() {
		m := scanner.Text()
		fmt.Print(m)
	}
	cmd.Wait()
}

func sync() {
	if _, err := os.Stat("dist"); os.IsNotExist(err) {
		log.Fatal("'dist' could not be found")
	}

	cmd := exec.Command("rsync",
		"-arvP",
		"-e", "ssh -i " + config.Ssh.Key,
		"dist",
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

func parseFlags()  {
	cliFlags = CliFlags{}

	// TODO: Use FlagSet to have both short, long options
	flag.BoolVar(&cliFlags.Build, "build", false, "Build project")
	flag.BoolVar(&cliFlags.Sync, "sync", false, "Sync files")
	flag.BoolVar(&cliFlags.Backup, "backup", false, "Backup target")
	flag.StringVar(&cliFlags.Env, "env", "", "Environment")

	flag.Parse()

	if cliFlags.Env == "" {
		log.Fatal("Missing '--env'")
	}
}

func loadConf() {
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

	log.Println(config)

	// make sure there is no trailing /
	config.Ssh.Path = strings.TrimSuffix(config.Ssh.Path, "/")
}

func main() {
	parseFlags()
	loadConf()

	projectDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Project Dir: %s\n", projectDir)

	if cliFlags.Build {
		build()
	}

	if cliFlags.Sync {
		sync()
	}
}