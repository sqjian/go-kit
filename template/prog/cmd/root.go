package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var (
	use = strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
)

var RootCmd = &cobra.Command{
	Use:   use,
	Short: fmt.Sprintf("%v is used to do xxx", use),
}

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
