// SPDX-License-Identifier: GPL-3.0-or-later
//
// Alpine Linux VM image management.
// Downloads and caches a minimal Alpine qcow2 image on first use.

package vm

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const alpineVersion = "3.19.1"

func alpineArch() string {
	if runtime.GOARCH == "arm64" {
		return "aarch64"
	}
	return "x86_64"
}

// alpineImageURL returns the download URL for the Alpine virtual disk image.
func alpineImageURL() string {
	return fmt.Sprintf(
		"https://dl-cdn.alpinelinux.org/alpine/v%s/releases/cloud/alpine-virt-%s-%s.qcow2",
		alpineVersion[:4], alpineVersion, alpineArch(),
	)
}

// cacheDir returns the OS-appropriate cache directory for linuxfs-mac.
func cacheDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("user cache dir: %w", err)
	}
	dir := filepath.Join(base, "linuxfs-mac")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("create cache dir: %w", err)
	}
	return dir, nil
}

// ensureAlpineImage returns the path to a cached Alpine qcow2 image,
// downloading it if not present.
func ensureAlpineImage(cfg Config) (string, error) {
	dir := cfg.DataDir
	if dir == "" {
		var err error
		dir, err = cacheDir()
		if err != nil {
			return "", err
		}
	}

	imageName := fmt.Sprintf("alpine-virt-%s-%s.qcow2", alpineVersion, alpineArch())
	imagePath := filepath.Join(dir, imageName)

	if _, err := os.Stat(imagePath); err == nil {
		return imagePath, nil
	}

	url := alpineImageURL()
	return "", fmt.Errorf(
		"Alpine VM image not found at %s\nDownload manually:\n  curl -L %s -o %s",
		imagePath, url, imagePath,
	)
}
