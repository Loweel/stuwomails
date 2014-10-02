package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

var appFiles = []string{"app/main.go", "app/assets.go"}

const configFileDestination = "assets/config.yaml"

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalln("No command specified")
		return
	}

	switch flag.Arg(0) {
	case "assets":
		buildAssets()
	case "build":
		buildBuild()
	case "get":
		buildGet()
	case "install":
		buildInstall()
	case "production":
		buildProduction()
	case "run":
		buildRun()
	default:
		log.Fatalf("Unknown command %q", flag.Arg(0))
	}
}

func buildAssets() {
	runPrint("go-bindata", "-o", "app/assets.go", "-prefix", "assets/", "assets/...")
	err := os.Remove(configFileDestination)

	if err != nil {
		log.Fatalf("Configuration file remove failed: %s", err)
	}
}

func buildBuild() {
	copyConfigFile("development")
	buildAssets()
	runSimpleGoCommand("build")
}

func buildGet() {
	runPrint("go", "get", "github.com/jteeuwen/go-bindata/...")
	copyConfigFile("development")
	buildAssets()
	runPrint("go", "get", "github.com/srhnsn/stuwomails/...")
}

func buildInstall() {
	cwd, err := os.Getwd()

	if err != nil {
		log.Fatalf("os.Getwd() failed: %s", err)
	}

	os.Setenv("GOBIN", cwd)
	buildAssets()
	runSimpleGoCommand("install")
}

func buildProduction() {
	os.Setenv("GOARCH", "amd64")
	os.Setenv("GOOS", "linux")

	copyConfigFile("production")
	buildInstall()
	os.Remove("stuwomails")
	os.Rename("main", "stuwomails")
}

func buildRun() {
	copyConfigFile("development")
	buildAssets()
	runSimpleGoCommand("run")
}

func copyConfigFile(name string) {
	dst, err := os.Create(configFileDestination)

	if err != nil {
		log.Fatalf("Cannot create configuration file: %s", err)
	}

	defer dst.Close()

	src, err := os.Open("resources/config." + name + ".yaml")

	if err != nil {
		log.Fatalf("Cannot read configuration file: %s", err)
	}

	defer src.Close()

	_, err = io.Copy(dst, src)

	if err != nil {
		log.Fatalf("Configuration file copy failed: %s", err)
	}
}

func runPrint(cmd string, args ...string) {
	log.Println(cmd, strings.Join(args, " "))

	ecmd := exec.Command(cmd, args...)

	ecmd.Stdout = os.Stdout
	ecmd.Stderr = os.Stderr

	err := ecmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func runSimpleGoCommand(name string) {
	runPrint("go", append([]string{name}, appFiles...)...)
}
