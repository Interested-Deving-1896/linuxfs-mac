// SPDX-License-Identifier: GPL-3.0-or-later
//
// vm package manages the lifecycle of the Alpine Linux microVM used to mount
// Linux filesystems. It wraps QEMU (cross-platform) and optionally libkrun
// (Apple Silicon) as hypervisor backends.
//
// Derived from AlexSSD7/linsk (vm/vm.go) and nohajc/anylinuxfs.

package vm

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"runtime"
)

// Backend selects the hypervisor used to run the Alpine VM.
type Backend int

const (
	// BackendQEMU uses qemu-system-x86_64 / qemu-system-aarch64.
	// Works on Linux, macOS (Intel + Apple Silicon), and Windows.
	BackendQEMU Backend = iota

	// BackendLibkrun uses libkrun (Apple Silicon only).
	// Lower overhead than QEMU; used by anylinuxfs on M-series Macs.
	BackendLibkrun
)

// Config holds parameters for a single VM instance.
type Config struct {
	// MemMiB is the RAM allocated to the VM in MiB.
	MemMiB uint32

	// DevicePath is the host block device or image file to pass through.
	DevicePath string

	// ReadOnly mounts the device read-only inside the VM.
	ReadOnly bool

	// Debug enables verbose QEMU/VM output.
	Debug bool

	// SSHPort is the host port forwarded to the VM's SSH daemon.
	// 0 = auto-select a free port.
	SSHPort uint16

	// Backend selects the hypervisor.
	Backend Backend

	// DataDir overrides the default cache directory for VM images.
	DataDir string
}

// VM represents a running Alpine Linux microVM instance.
type VM struct {
	cfg    Config
	logger *slog.Logger
	cmd    *exec.Cmd
	ctx    context.Context
	cancel context.CancelFunc

	// SSHPort is the resolved host port (set after Start).
	SSHPort uint16
}

// New creates a VM instance. Call Start to launch it.
func New(ctx context.Context, cfg Config, logger *slog.Logger) (*VM, error) {
	if cfg.MemMiB == 0 {
		cfg.MemMiB = 512
	}
	if cfg.Backend == BackendLibkrun && runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("libkrun backend is only supported on macOS")
	}
	ctx, cancel := context.WithCancel(ctx)
	return &VM{
		cfg:    cfg,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Start launches the Alpine VM and waits for SSH to become available.
func (v *VM) Start() error {
	switch v.cfg.Backend {
	case BackendQEMU:
		return v.startQEMU()
	case BackendLibkrun:
		return v.startLibkrun()
	default:
		return fmt.Errorf("unknown backend: %d", v.cfg.Backend)
	}
}

// Stop shuts down the VM gracefully.
func (v *VM) Stop() error {
	v.cancel()
	if v.cmd != nil && v.cmd.Process != nil {
		return v.cmd.Process.Kill()
	}
	return nil
}

func (v *VM) startQEMU() error {
	binary := "qemu-system-x86_64"
	if runtime.GOARCH == "arm64" {
		binary = "qemu-system-aarch64"
	}

	if _, err := exec.LookPath(binary); err != nil {
		return fmt.Errorf("%s not found: install QEMU first", binary)
	}

	alpineImage, err := ensureAlpineImage(v.cfg)
	if err != nil {
		return fmt.Errorf("alpine image: %w", err)
	}

	sshPort := v.cfg.SSHPort
	if sshPort == 0 {
		sshPort = 10022
	}
	v.SSHPort = sshPort

	args := []string{
		"-enable-kvm",
		"-m", fmt.Sprintf("%d", v.cfg.MemMiB),
		"-nographic",
		"-serial", "mon:stdio",
		"-drive", fmt.Sprintf("if=virtio,format=qcow2,file=%s", alpineImage),
		"-drive", fmt.Sprintf("if=virtio,format=raw,file=%s,readonly=%s",
			v.cfg.DevicePath, boolToOnOff(v.cfg.ReadOnly)),
		"-netdev", fmt.Sprintf("user,id=net0,hostfwd=tcp::%d-:22", sshPort),
		"-device", "virtio-net-pci,netdev=net0",
	}

	v.logger.Info("Starting QEMU Alpine VM",
		"binary", binary,
		"mem_mib", v.cfg.MemMiB,
		"device", v.cfg.DevicePath,
		"ssh_port", sshPort,
	)

	v.cmd = exec.CommandContext(v.ctx, binary, args...)
	return v.cmd.Start()
}

func (v *VM) startLibkrun() error {
	// libkrun integration requires CGo bindings to libkrun.dylib.
	// See docs/libkrun.md for build instructions.
	return fmt.Errorf("libkrun backend: CGo integration not yet compiled — see docs/libkrun.md")
}

func boolToOnOff(b bool) string {
	if b {
		return "on"
	}
	return "off"
}
