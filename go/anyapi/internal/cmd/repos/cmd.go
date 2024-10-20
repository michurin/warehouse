package repos

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"anyapi/internal/api"
	"anyapi/internal/cmd"
	"anyapi/internal/str"
)

func Cmd(root cmd.CmdInterface) {
	var verbose bool

	cmdModerate := &cobra.Command{
		Use:   "repos [username]",
		Short: "Show repos of user",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, argsv []string) error {
			for _, user := range argsv {
				resp, err := api.GetStruct("http://api.github.com/users/"+url.PathEscape(user)+"/repos", verbose)
				cobra.CheckErr(err)
				tweak(resp)
				fmt.Println(str.Repr(resp))
			}
			return nil
		},
	}

	cmdModerate.Flags().BoolVarP(&verbose, "verbose", "v", false, "be verbose")

	root.AddCommand(cmdModerate)
}

func tweak(x any) {
	switch a := x.(type) {
	case []any:
		for _, v := range a {
			tweak(v)
		}
	case map[string]any:
		for k, v := range a {
			_ = v
			if k == "watchers" {
				n, ok := v.(float64)
				if ok {
					a[k] = str.ValueWithComment{Value: str.Repr(v), Comment: strings.Repeat("★", min(int(n), 50))} //nolint:mnd
				}
			}
		}
	}
}
