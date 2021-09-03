package cmd

import "github.com/spf13/cobra"

var binCmd3 = &cobra.Command{
	Use:   "cmd3",
	Short: "cmd3 message",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	RootCmd.AddCommand(binCmd3)
}
