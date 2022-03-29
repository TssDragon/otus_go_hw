package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (exitCode int) {
	argCmd := cmd[0]
	args := cmd[1:]
	localEnv := makeEnvAsStringSlice(env)

	return realRunCmd(argCmd, args, localEnv)
}

func makeEnvAsStringSlice(env Environment) []string {
	sliceEnv := []string{}
	for key, value := range env {
		_, ok := os.LookupEnv(key)
		if ok {
			os.Unsetenv(key)
		}

		if !value.NeedRemove {
			sliceEnv = append(sliceEnv, key+"="+value.Value)
		}
	}
	return sliceEnv
}

func realRunCmd(cmd string, args []string, env []string) (exitCode int) {
	execCmd := exec.Command(cmd, args...)

	execCmd.Env = append(os.Environ(), env...)
	execCmd.Stdout = os.Stdout
	execCmd.Stdin = os.Stdin
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		fmt.Println(err)
		state := execCmd.ProcessState
		return state.ExitCode()
	}
	return
}
