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

const (
	alpineVersion   = "3.19.1"
	alpineArch_amd  = "x86_64"
	alpineArch_arm  = "aarch64"
	alpineImageSize = "512M"
)

// alpineImageURL returns the download URL for the Alpine virtual disk image.
func alpineImageURL() string {
	arch := alpineArch_amd
	if runtime.GOARCH == "arm64" {
		arch = alpineArch_arm
	}
	return fmt.Sprintf(
		"https://dl-cdn.alpinelinux.org/alpine/v%s/releases/cloud/alpine-virt-%s-%s.qcow2",
		alpineVersion[:4], alpineVersion, arch,
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
	dir := cfg.DevicePath // reuse data-dir if set; otherwise use cache
	if dir == "" {
		var err error
		dir, err = cacheDir()
		if err != nil {
			return "", err
		}
	}

	arch := alpineArch_amd
	if runtime.GOARCH == "arm64" {
		arch = alpineArch_arm
	}
	imageName := fmt.Sprintf("alpine-virt-%s-%s.qcow2", alpineVersion, arch)
	imagePath := filepath.Join(dir, imageName)

	if _, err := os.Stat(imagePath); err == nil {
		return imagePath, nil // already cached
	}

	url := alpineImageURL()
	fmt.Printf("Downloading Alpine Linux VM image (%s) ...\n", alpineVersion)
	fmt.Printf("  URL: %s\n", url)
	fmt.Printf("  Destination: %s\n", imagePath)

	// TODO: implement download with progress bar.
	// For now, return an error directing the user to download manually.
	return "", fmt.Errorf(
		"Alpine VM image not found at %s\n"+
			"Download manually:\n  curl -L %s -o %s",
		imagePath, url, imagePath,
	)
}
