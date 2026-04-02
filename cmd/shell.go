// SPDX-License-Identifier: GPL-3.0-or-later

package cmd

import (
	"flag"
	"fmt"
	"os"
)

func runShell(args []string) {
	fs := flag.NewFlagSet("shell", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}
	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Usage: linuxfs-mac shell <device>")
		os.Exit(1)
	}
	fmt.Printf("Starting Alpine VM for %s ...\n", fs.Arg(0))
	// TODO: start VM, wait for SSH, exec interactive session.
	fmt.Fprintln(os.Stderr, "shell: not yet implemented")
	os.Exit(1)
}
