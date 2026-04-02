// SPDX-License-Identifier: GPL-3.0-or-later

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("linuxfs-mac %s\n", Version)
		fmt.Println("Upstream projects: AlexSSD7/linsk, nohajc/anylinuxfs")
	},
}
