// SPDX-License-Identifier: GPL-3.0-or-later

package cmd

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

func runMount(args []string) {
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	mountPoint  := fs.String("mountpoint", "", "Host mount point (default: auto-generated)")
	readOnly    := fs.Bool("read-only", false, "Mount read-only")
	luks        := fs.String("luks", "", "In-VM device for LUKS container (e.g. /dev/sda1)")
	lvm         := fs.String("lvm", "", "LVM volume group/logical volume (e.g. vg0/home)")
	fstype      := fs.String("fstype", "", "Filesystem type hint (ext4, btrfs, xfs, zfs, ...)")
	mountOpts   := fs.String("mount-opts", "", "Extra mount options passed inside the VM")
	netShare    := fs.Bool("network-share", false, "Expose share on the local network")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: linuxfs-mac mount [flags] <device>

Mount a block device or disk image containing a Linux filesystem.
The device is passed through to an Alpine Linux VM which mounts it natively.
The mounted filesystem is then exposed to the host via a network share.

Examples:
  linuxfs-mac mount /dev/disk2s1
  linuxfs-mac mount /dev/disk2s1 --luks /dev/sda1
  linuxfs-mac mount /dev/disk2s1 --lvm vg0/home
  linuxfs-mac mount /dev/disk2s1 --read-only
  linuxfs-mac mount /path/to/disk.img
  linuxfs-mac mount /dev/disk2s1 --fstype btrfs --mount-opts subvol=@home

Flags:
`)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}
	if fs.NArg() < 1 {
		fs.Usage()
		os.Exit(1)
	}

	device := fs.Arg(0)
	backend := flagBackend
	if backend == "" {
		backend = defaultBackend()
	}

	fmt.Printf("Mounting %s via %s backend ...\n", device, backend)
	fmt.Printf("  VM memory:  %d MiB\n", flagVMMemMiB)
	fmt.Printf("  Read-only:  %v\n", *readOnly)
	if *luks != "" {
		fmt.Printf("  LUKS:       %s\n", *luks)
	}
	if *lvm != "" {
		fmt.Printf("  LVM:        %s\n", *lvm)
	}
	if *fstype != "" {
		fmt.Printf("  FS type:    %s\n", *fstype)
	}
	if *mountPoint != "" {
		fmt.Printf("  Mountpoint: %s\n", *mountPoint)
	}
	if *mountOpts != "" {
		fmt.Printf("  Mount opts: %s\n", *mountOpts)
	}
	if *netShare {
		fmt.Println("  Network share: enabled")
	}

	// TODO: wire up vm.Manager to start Alpine VM, pass device through,
	// mount inside VM, and start share backend.
	fmt.Fprintln(os.Stderr, "mount: VM backend not yet wired — see vm/ package")
	os.Exit(1)
}

func runUnmount(args []string) {
	fs := flag.NewFlagSet("unmount", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}
	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Usage: linuxfs-mac unmount <device|mountpoint>")
		os.Exit(1)
	}
	fmt.Printf("Unmounting %s ...\n", fs.Arg(0))
	// TODO: signal running VM to unmount and shut down.
	fmt.Fprintln(os.Stderr, "unmount: not yet implemented")
	os.Exit(1)
}

func defaultBackend() string {
	switch runtime.GOOS {
	case "darwin":
		return "afp"
	case "windows":
		return "smb"
	default:
		return "ftp"
	}
}
