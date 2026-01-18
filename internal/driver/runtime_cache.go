package driver

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	coruntime "github.com/0xmukesh/coco/internal/runtime"
)

type RuntimeCache struct {
	dir     string
	objPath string
}

func NewRuntimeCache() (*RuntimeCache, error) {
	var cacheDir string
	switch runtime.GOOS {
	case "linux", "darwin":
		cacheDir = filepath.Join(os.Getenv("HOME"), ".cache", "coco")
	case "windows":
		cacheDir = filepath.Join(os.Getenv("LOCALAPPDATA"), "coco", "cache")
	default:
		cacheDir = filepath.Join(os.TempDir(), "coco-cache")
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	runtimeObjPath := filepath.Join(cacheDir, "runtime.o")

	cache := &RuntimeCache{
		dir:     cacheDir,
		objPath: runtimeObjPath,
	}

	if _, err := os.Stat(runtimeObjPath); os.IsNotExist(err) {
		if err := cache.compileRuntimeFile(); err != nil {
			return nil, fmt.Errorf("failed to compile runtime: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to check runtime file: %w", err)
	}

	return cache, nil
}

func (rc *RuntimeCache) compileRuntimeFile() error {
	runtimeSource := coruntime.Source
	tmpSrc := filepath.Join(rc.dir, "runtime.c")

	if err := os.WriteFile(tmpSrc, runtimeSource, 0644); err != nil {
		return fmt.Errorf("failed to write runtime source: %w", err)
	}
	defer os.Remove(tmpSrc)

	cmd := exec.Command("clang", "-c", "-O1", "-fPIC", tmpSrc, "-o", rc.objPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clang compilation failed: %w", err)
	}

	fmt.Printf("compiled runtime library (coco_runtime.c) to %s\n", rc.objPath)
	return nil
}
