package cmd

import "github.com/spf13/cobra"

var binCmd2 = &cobra.Command{
	Use:     "cmd2",
	Short:   "cmd2 message",
	GroupID: "2",
	Run:     func(cmd *cobra.Command, args []string) {},
}

func init() {
	RootCmd.AddCommand(binCmd2)
}
