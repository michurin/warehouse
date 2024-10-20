package cmd

import "github.com/spf13/cobra"

type CmdInterface interface {
	AddCommand(cmds ...*cobra.Command)
}
