package cmd

import (
	"github.com/spf13/cobra"
)

var binCmd1 = &cobra.Command{
	Use:   "cmd1",
	Short: "cmd1 message",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	RootCmd.AddCommand(binCmd1)
}
