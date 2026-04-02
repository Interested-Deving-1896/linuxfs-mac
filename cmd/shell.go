// SPDX-License-Identifier: GPL-3.0-or-later

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell <device>",
	Short: "Open an interactive shell inside the Alpine VM for a device.",
	Long: `Start the Alpine VM with the given device attached and drop into an
interactive SSH shell. Useful for manual inspection, fsck, or LVM/LUKS setup.

Example:
  linuxfs-mac shell /dev/disk2s1`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Starting Alpine VM for %s ...\n", args[0])
		// TODO: start VM, wait for SSH, exec interactive session.
		fmt.Fprintln(os.Stderr, "shell: not yet implemented")
		return nil
	},
}
