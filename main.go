// SPDX-License-Identifier: GPL-3.0-or-later
//
// linuxfs-mac — access Linux-native filesystems on macOS and Windows.
//
// Unified from:
//   - AlexSSD7/linsk  (Go core, QEMU VM, SMB/AFP/FTP backends)
//   - nohajc/anylinuxfs (libkrun microVM, NFS backend, Apple Silicon)
//
// Architecture: spins up a lightweight Alpine Linux VM, passes through the
// target block device, mounts it inside the VM, then exposes it to the host
// via a network share (SMB on Windows, AFP/NFS on macOS, FTP as fallback).

package main

import (
	"github.com/linuxfs-mac/linuxfs-mac/cmd"
)

func main() {
	cmd.Execute()
}
