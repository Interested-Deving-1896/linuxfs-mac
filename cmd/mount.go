// SPDX-License-Identifier: GPL-3.0-or-later

package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	flagMountPoint  string
	flagReadOnly    bool
	flagLUKSDevice  string
	flagLVMGroup    string
	flagFSType      string
	flagMountOpts   string
	flagNetworkShare bool
)

var mountCmd = &cobra.Command{
	Use:   "mount <device>",
	Short: "Mount a Linux filesystem and expose it as a network share.",
	Long: `Mount a block device or disk image containing a Linux filesystem.

The device is passed through to an Alpine Linux VM which mounts it natively.
The mounted filesystem is then exposed to the host via a network share.

Examples:
  # Mount an ext4 partition
  linuxfs-mac mount /dev/disk2s1

  # Mount a LUKS-encrypted partition
  linuxfs-mac mount /dev/disk2s1 --luks /dev/sda1

  # Mount a specific LVM logical volume
  linuxfs-mac mount /dev/disk2s1 --lvm my-vg/my-lv

  # Mount read-only
  linuxfs-mac mount /dev/disk2s1 --read-only

  # Mount a disk image
  linuxfs-mac mount /path/to/disk.img

  # Mount with explicit filesystem type
  linuxfs-mac mount /dev/disk2s1 --fstype btrfs --mount-opts subvol=@home

  # Share on the local network (not just localhost)
  linuxfs-mac mount /dev/disk2s1 --network-share`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		device := args[0]

		backend := flagBackend
		if backend == "" {
			backend = defaultBackend()
		}

		fmt.Printf("Mounting %s via %s backend ...\n", device, backend)
		fmt.Printf("  VM memory:  %d MiB\n", flagVMMemMiB)
		fmt.Printf("  Read-only:  %v\n", flagReadOnly)
		if flagLUKSDevice != "" {
			fmt.Printf("  LUKS:       %s\n", flagLUKSDevice)
		}
		if flagLVMGroup != "" {
			fmt.Printf("  LVM:        %s\n", flagLVMGroup)
		}
		if flagFSType != "" {
			fmt.Printf("  FS type:    %s\n", flagFSType)
		}

		// TODO: wire up vm.Manager to start Alpine VM, pass device through,
		// mount inside VM, and start share backend.
		fmt.Fprintln(os.Stderr, "mount: VM backend not yet wired — see vm/ package")
		return nil
	},
}

var unmountCmd = &cobra.Command{
	Use:   "unmount <device|mountpoint>",
	Short: "Unmount a previously mounted Linux filesystem.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Unmounting %s ...\n", args[0])
		// TODO: signal running VM to unmount and shut down.
		fmt.Fprintln(os.Stderr, "unmount: not yet implemented")
		return nil
	},
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

func init() {
	mountCmd.Flags().StringVar(&flagMountPoint, "mountpoint", "",
		"Host mount point (default: auto-generated under /Volumes on macOS)")
	mountCmd.Flags().BoolVar(&flagReadOnly, "read-only", false,
		"Mount read-only")
	mountCmd.Flags().StringVar(&flagLUKSDevice, "luks", "",
		"In-VM device name for LUKS container (e.g. /dev/sda1)")
	mountCmd.Flags().StringVar(&flagLVMGroup, "lvm", "",
		"LVM volume group/logical volume to mount (e.g. vg0/home)")
	mountCmd.Flags().StringVar(&flagFSType, "fstype", "",
		"Filesystem type hint (ext4, btrfs, xfs, zfs, ntfs, exfat, ...)")
	mountCmd.Flags().StringVar(&flagMountOpts, "mount-opts", "",
		"Extra options passed to mount inside the VM (e.g. subvol=@home)")
	mountCmd.Flags().BoolVar(&flagNetworkShare, "network-share", false,
		"Expose share on the local network, not just localhost")
}
