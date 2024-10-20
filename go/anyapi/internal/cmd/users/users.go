package users

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"anyapi/internal/api"
	"anyapi/internal/cmd"
	"anyapi/internal/str"
)

func Cmd(root cmd.CmdInterface) {
	var verbose bool

	cmdModerate := &cobra.Command{
		Use:   "users [username]",
		Short: "Show user information",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, argsv []string) error {
			for _, user := range argsv {
				resp, err := api.GetStruct("http://api.github.com/users/"+url.PathEscape(user), verbose)
				cobra.CheckErr(err)
				fmt.Println(str.Repr(resp))
			}
			return nil
		},
	}

	cmdModerate.Flags().BoolVarP(&verbose, "verbose", "v", false, "be verbose")

	root.AddCommand(cmdModerate)
}
