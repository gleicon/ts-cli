package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func help() {
	dotdir := filepath.Join(GetEnvHomeDir(), DOTDIR)
	fmt.Println("ts-cli -> Time Series for Command LIne")
	fmt.Println("Stores and chart a timeseries for any label and value on terminal")
	fmt.Println("Data stored in BoltDB at ~/.ts-cli/timeseries.db")
	fmt.Println("echo 'mymetric=10' | ts-cli to store a new point using the current time. data format is \"label=value\"")
	fmt.Println("ts-cli -p <mymetric> to print the points for the label mymetric")
	fmt.Println("Data dir:", dotdir)
	os.Exit(0)
}

func GetEnvHomeDir() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	}
	return os.Getenv(env)
}
