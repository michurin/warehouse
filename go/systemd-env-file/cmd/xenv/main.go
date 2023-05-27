package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/michurin/systemd-env-file/sdenv"
)

func app(args []string, stdout, stderr io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("you are to specify command")
	}
	err := (error)(nil)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env, err = sdenv.Environ(os.Environ(), "xenv.env")
	if err != nil {
		return err
	}
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := app(os.Args[1:], os.Stdout, os.Stderr)
	if err != nil {
		fmt.Println(err)
	}
}
