// SPDX-License-Identifier: GPL-3.0-or-later
//
// mount/backends.go — file share backend registry.
//
// Each backend exposes the filesystem mounted inside the Alpine VM to the host.
// Derived from AlexSSD7/linsk (share/) and nohajc/anylinuxfs (NFS backend).

package mount

import (
	"fmt"
	"runtime"
)

// Backend identifies a network file share protocol.
type Backend string

const (
	BackendAFP Backend = "afp"  // Apple Filing Protocol — macOS default
	BackendNFS Backend = "nfs"  // NFS — macOS alternative (anylinuxfs default)
	BackendSMB Backend = "smb"  // SMB/CIFS — Windows default
	BackendFTP Backend = "ftp"  // FTP — cross-platform fallback
)

// DefaultBackend returns the recommended backend for the current OS.
func DefaultBackend() Backend {
	switch runtime.GOOS {
	case "darwin":
		return BackendAFP
	case "windows":
		return BackendSMB
	default:
		return BackendFTP
	}
}

// Config holds share server configuration.
type Config struct {
	Backend    Backend
	ListenIP   string
	ListenPort uint16
	// NetworkShare allows connections from outside localhost.
	NetworkShare bool
}

// ShareInfo describes a running share that the host can connect to.
type ShareInfo struct {
	Backend    Backend
	MountURL   string // e.g. afp://127.0.0.1/linuxfs or smb://127.0.0.1/linuxfs
	MountPoint string // where it was auto-mounted on the host
}

// Validate checks that the backend is supported on the current OS.
func (c Config) Validate() error {
	switch c.Backend {
	case BackendAFP:
		if runtime.GOOS != "darwin" {
			return fmt.Errorf("AFP backend is only supported on macOS")
		}
	case BackendNFS:
		if runtime.GOOS != "darwin" {
			return fmt.Errorf("NFS backend is currently only supported on macOS")
		}
	case BackendSMB, BackendFTP:
		// cross-platform
	default:
		return fmt.Errorf("unknown backend: %q", c.Backend)
	}
	return nil
}

// MountURL builds the URL the host uses to connect to the share.
func (c Config) MountURL(shareName string) string {
	ip := c.ListenIP
	if ip == "" {
		ip = "127.0.0.1"
	}
	switch c.Backend {
	case BackendAFP:
		return fmt.Sprintf("afp://%s/%s", ip, shareName)
	case BackendNFS:
		return fmt.Sprintf("nfs://%s/%s", ip, shareName)
	case BackendSMB:
		return fmt.Sprintf("smb://%s/%s", ip, shareName)
	case BackendFTP:
		port := c.ListenPort
		if port == 0 {
			port = 2121
		}
		return fmt.Sprintf("ftp://%s:%d", ip, port)
	default:
		return ""
	}
}
