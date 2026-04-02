// SPDX-License-Identifier: GPL-3.0-or-later

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List currently mounted Linux filesystems.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: read state from flagDataDir / OS cache dir and print active mounts.
		fmt.Println("No active mounts.")
		return nil
	},
}
