package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Fatalf("Not enough arguments")
	}

	envDir := args[1]
	env, err := ReadDir(envDir)
	if err != nil {
		log.Fatalf("ReadDir err: %v", err)
	}
	args = args[2:]

	exitCode := RunCmd(args, env)
	os.Exit(exitCode)
}
