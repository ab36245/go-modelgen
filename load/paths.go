package load

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func pathsToFiles(files *[]string, paths []string) error {
	for _, path := range paths {
		if err := pathToFiles(files, path); err != nil {
			return err
		}
	}
	return nil
}

func pathToFiles(files *[]string, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		*files = append(*files, path)
		return nil
	}
	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), defExtension) {
			return nil
		}
		*files = append(*files, path)
		return nil
	})
	return err
}
