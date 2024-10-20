package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"anyapi/internal/cmd"
	"anyapi/internal/cmd/repos"
	"anyapi/internal/cmd/users"
)

func main() {
	rootCmd := &cobra.Command{Use: "anyapi"}
	for _, b := range []func(cmd.CmdInterface){
		users.Cmd,
		repos.Cmd,
	} {
		b(rootCmd)
	}
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
	}
}
